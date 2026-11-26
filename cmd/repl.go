package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"slices"
	"strconv"
	"strings"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

const (
	replPrompt = "\x1b[32mha > \x1b[0m"

	// replLoginExitCode asks the caller of the REPL to start a login shell
	replLoginExitCode = 10

	replClearScreen = "\x1b[2J\x1b[H"

	keyCtrlC = 3
	keyCtrlD = 4
)

// replBuiltins are the commands of the REPL itself, with their descriptions
var replBuiltins = [][2]string{
	{"exit", "Exit the interactive prompt"},
	{"login", "Leave the interactive prompt for the operating system shell"},
}

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Starts an interactive prompt to run the CLI commands from",
	Long: `
Starts an interactive prompt with line editing, history and Tab completion.
The commands are entered without the leading 'ha'.

On top of the CLI commands the prompt understands:

  exit     Exit the prompt
  login    Exit the prompt with the exit code 10, the caller is expected to
           start a login shell then

Ctrl+C cancels the running command or the line being edited. Ctrl+D on an
empty line resets the prompt to its initial state.`,
	Example: `
  ha repl
  ha repl --no-banner`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("repl", "args", args)

		stdin := int(os.Stdin.Fd())
		if !term.IsTerminal(stdin) {
			helper.PrintErrorString("the interactive prompt needs a terminal")
			ExitWithError = true
			return
		}

		self, err := os.Executable()
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
			return
		}

		// Run the commands with the global flags the REPL was started with
		var globalArgs []string
		cmd.Root().PersistentFlags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				globalArgs = append(globalArgs, "--"+f.Name+"="+f.Value.String())
			}
		})

		// The arguments of the command come from the user of the prompt
		selfCommand := func(ctx context.Context, args []string) *exec.Cmd {
			return exec.CommandContext(ctx, self, append(slices.Clip(globalArgs), args...)...) //nolint:gosec
		}

		noBanner, _ := cmd.Flags().GetBool("no-banner")

		r := &repl{
			in:     os.Stdin,
			out:    os.Stdout,
			banner: !noBanner,
			run: func(args []string) {
				c := selfCommand(context.Background(), args)
				c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
				err := c.Run()
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) {
					if err != nil {
						helper.PrintError(err)
					}
					return
				}
				// The terminal echoes ^C without a line break
				if exitErr.ExitCode() == -1 {
					fmt.Println()
				}
			},
			complete: func(ctx context.Context, args []string) string {
				args = append([]string{cobra.ShellCompRequestCmd}, args...)
				out, _ := selfCommand(ctx, args).Output()
				return string(out)
			},
			raw: func() (func(), error) {
				state, err := term.MakeRaw(stdin)
				if err != nil {
					return nil, err
				}
				return func() { _ = term.Restore(stdin, state) }, nil
			},
			size: func() (int, int, error) { return term.GetSize(stdin) },
		}

		ExitCode = r.loop()
	},
}

func init() {
	replCmd.Flags().Bool("no-banner", false, "Don't print the banner when the prompt starts")
	replCmd.Flags().Lookup("no-banner").NoOptDefVal = "true"
	replCmd.RegisterFlagCompletionFunc("no-banner", boolCompletions)

	rootCmd.AddCommand(replCmd)
}

// repl is the interactive prompt. The terminal specifics are behind function
// fields, so the tests can drive it without a terminal.
type repl struct {
	in     io.Reader
	out    io.Writer
	banner bool
	// run executes the CLI with args and waits for it to finish
	run func(args []string)
	// complete returns the output of the cobra completion request for args.
	// It must return early when ctx is cancelled.
	complete func(ctx context.Context, args []string) string
	// raw puts the terminal into raw mode and returns the way back
	raw func() (restore func(), err error)
	// size returns the width and height of the terminal
	size func() (width, height int, err error)
}

// keyReader reads the keys for term.Terminal one byte at a time. It must
// not read ahead, because input following a command line belongs to that
// command, e.g. a pasted password.
type keyReader struct {
	in io.Reader
	// key is the last key read, if it was Ctrl+C or Ctrl+D. term.Terminal
	// reports both as io.EOF, the same as the end of the input.
	key byte
	// pending delivers the result of a read which is under way
	pending chan keyRead
	// backlog holds the keys read by waitInterrupt
	backlog []keyRead
}

type keyRead struct {
	key byte
	n   int
	err error
}

// next starts reading a key, unless a read is under way already
func (k *keyReader) next() <-chan keyRead {
	if k.pending == nil {
		pending := make(chan keyRead, 1)
		go func() {
			var buf [1]byte
			n, err := k.in.Read(buf[:])
			pending <- keyRead{buf[0], n, err}
		}()
		k.pending = pending
	}
	return k.pending
}

func (k *keyReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	var read keyRead
	if len(k.backlog) > 0 {
		read, k.backlog = k.backlog[0], k.backlog[1:]
	} else {
		read = <-k.next()
		k.pending = nil
	}
	k.key = 0
	if read.n > 0 {
		p[0] = read.key
		if read.key == keyCtrlC || read.key == keyCtrlD {
			k.key = read.key
		}
	}
	return read.n, read.err
}

// waitInterrupt waits for done to be closed and reports if Ctrl+C was
// pressed before that. The keys read meanwhile are kept for Read, including
// the Ctrl+C. Call it only while term.Terminal waits for a callback, as its
// next Read has to pick up the read left under way.
func (k *keyReader) waitInterrupt(done <-chan struct{}) bool {
	for {
		select {
		case <-done:
			return false
		case read := <-k.next():
			k.pending = nil
			k.backlog = append(k.backlog, read)
			if read.err != nil {
				<-done
				return false
			}
			if read.n > 0 && read.key == keyCtrlC {
				return true
			}
		}
	}
}

// replHistory keeps the entered lines for term.Terminal. Unlike the history
// of the terminal it skips empty lines and immediate repeats.
type replHistory struct {
	lines []string
}

func (h *replHistory) Add(entry string) {
	entry = strings.TrimSpace(entry)
	if entry == "" || (len(h.lines) > 0 && h.lines[len(h.lines)-1] == entry) {
		return
	}
	h.lines = append(h.lines, entry)
}

func (h *replHistory) Len() int {
	return len(h.lines)
}

func (h *replHistory) At(idx int) string {
	return h.lines[len(h.lines)-1-idx]
}

// loop reads and runs commands until the prompt is left and returns the exit
// code of the CLI
func (r *repl) loop() int {
	// Ctrl+C is meant for the running command, which shares the terminal
	// with the REPL. The signal must be caught and not ignored, as a command
	// would inherit the ignored signal.
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	defer signal.Stop(interrupts)

	keys := &keyReader{in: r.in}
	newTerminal := func(history term.History) *term.Terminal {
		t := term.NewTerminal(struct {
			io.Reader
			io.Writer
		}{keys, r.out}, replPrompt)
		t.History = history
		t.AutoCompleteCallback = func(line string, pos int, key rune) (string, int, bool) {
			if key != '\t' {
				return "", 0, false
			}
			// A completion asking Supervisor can take long, let Ctrl+C end it
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var newLine, list string
			var newPos int
			done := make(chan struct{})
			go func() {
				defer close(done)
				newLine, newPos, list = r.completeLine(ctx, line, pos)
			}()
			if keys.waitInterrupt(done) {
				cancel()
				<-done
				return line, pos, true
			}
			if list != "" {
				// Keep the completed line above the candidates
				_, _ = t.Write([]byte(replPrompt + line + "\n" + list))
			}
			return newLine, newPos, true
		}
		return t
	}

	if r.banner {
		r.run([]string{"banner"})
	}
	t := newTerminal(&replHistory{})

	for {
		line, err := r.readLine(t)
		switch {
		case err == nil:
		case errors.Is(err, io.EOF) && keys.key == keyCtrlC:
			// The terminal is replaced to drop the line being edited
			fmt.Fprintln(r.out, "^C")
			t = newTerminal(t.History)
			continue
		case errors.Is(err, io.EOF) && keys.key == keyCtrlD:
			fmt.Fprint(r.out, replClearScreen)
			if r.banner {
				r.run([]string{"banner"})
			}
			t = newTerminal(&replHistory{})
			continue
		case errors.Is(err, io.EOF):
			fmt.Fprintln(r.out)
			return 0
		default:
			helper.PrintError(err)
			return 1
		}

		args, err := splitLine(line)
		if err != nil {
			helper.PrintError(err)
			continue
		}
		if len(args) > 0 && args[0] == "ha" {
			fmt.Fprintln(r.out, "Note: Leading 'ha' is not necessary in this HA CLI")
			args = args[1:]
		}
		if len(args) == 0 {
			continue
		}

		switch {
		case args[0] == "exit" && len(args) == 1:
			return 0
		case args[0] == "login":
			return replLoginExitCode
		case args[0] == "help" && len(args) == 1:
			fmt.Fprintln(r.out, "Note: Use \"login\" to enter operating system shell")
			fmt.Fprintln(r.out)
			fmt.Fprintln(r.out, "--------------------------------")
		}
		r.run(args)
	}
}

// readLine reads one line with the terminal in raw mode. The mode is left
// again before returning, as the commands expect an ordinary terminal.
func (r *repl) readLine(t *term.Terminal) (string, error) {
	if r.raw != nil {
		restore, err := r.raw()
		if err != nil {
			return "", err
		}
		defer restore()
	}
	if r.size != nil {
		// A terminal of unknown size reports zero, keep the default then
		if width, height, err := r.size(); err == nil && width > 0 {
			_ = t.SetSize(width, height)
		}
	}
	return t.ReadLine()
}

// completeLine completes the word in front of pos. It returns the new line
// and cursor position, and the candidates to print when there are several
// and the word cannot be extended any further.
func (r *repl) completeLine(ctx context.Context, line string, pos int) (newLine string, newPos int, list string) {
	head := line[:pos]
	// Quoting is not understood here, leave such lines alone
	if strings.ContainsAny(head, `'"\`) {
		return line, pos, ""
	}

	words := strings.Fields(head)
	if len(words) == 0 || strings.TrimRight(head, " \t") != head {
		words = append(words, "")
	}
	if len(words) > 1 && words[0] == "ha" {
		words = words[1:]
	}

	candidates, noSpace := parseCompletions(r.complete(ctx, words))
	if len(words) == 1 {
		candidates = append(candidates, replBuiltins...)
	}

	// Cobra returns the values of --flag=value without the flag
	word := words[len(words)-1]
	if flag, value, ok := strings.Cut(word, "="); ok && strings.HasPrefix(flag, "-") {
		word = value
	}

	var matches [][2]string
	for _, candidate := range candidates {
		if strings.HasPrefix(candidate[0], word) {
			matches = append(matches, candidate)
		}
	}
	if len(matches) == 0 {
		return line, pos, ""
	}

	common := matches[0][0]
	width := 0
	for _, match := range matches {
		for !strings.HasPrefix(match[0], common) {
			common = common[:len(common)-1]
		}
		width = max(width, len(match[0]))
	}

	if len(matches) > 1 && common == word {
		var b strings.Builder
		for _, match := range matches {
			fmt.Fprintf(&b, "%-*s  %s\n", width, match[0], match[1])
		}
		return line, pos, b.String()
	}

	if len(matches) == 1 && !noSpace {
		common += " "
	}
	head = head[:len(head)-len(word)] + common
	return head + line[pos:], len(head), ""
}

// parseCompletions parses the output of a cobra completion request into
// pairs of a candidate and its description. noSpace reports if no space
// should follow the completed word.
func parseCompletions(output string) (candidates [][2]string, noSpace bool) {
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	last := lines[len(lines)-1]
	if !strings.HasPrefix(last, ":") {
		return nil, false
	}
	number, err := strconv.Atoi(last[1:])
	directive := cobra.ShellCompDirective(number)
	// File names are not completed. With these directives the lines are
	// file extensions or directories to look into, not candidates.
	const noCandidates = cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveFilterFileExt | cobra.ShellCompDirectiveFilterDirs
	if err != nil || directive&noCandidates != 0 {
		return nil, false
	}

	for _, line := range lines[:len(lines)-1] {
		// Messages for the user are not candidates
		if line == "" || strings.HasPrefix(line, "_activeHelp_ ") {
			continue
		}
		name, description, _ := strings.Cut(line, "\t")
		candidates = append(candidates, [2]string{name, description})
	}
	return candidates, directive&cobra.ShellCompDirectiveNoSpace != 0
}

// splitLine splits a line into arguments. Whitespace separates them, unless
// it is quoted with single quotes, double quotes or a backslash.
func splitLine(line string) ([]string, error) {
	var args []string
	var arg strings.Builder
	inArg := false
	var quote rune
	escaped := false

	for _, c := range line {
		switch {
		case escaped:
			arg.WriteRune(c)
			escaped = false
		case c == '\\' && quote != '\'':
			escaped = true
			inArg = true
		case quote != 0:
			if c == quote {
				quote = 0
			} else {
				arg.WriteRune(c)
			}
		case c == '\'' || c == '"':
			quote = c
			inArg = true
		case c == ' ' || c == '\t':
			if inArg {
				args = append(args, arg.String())
				arg.Reset()
				inArg = false
			}
		default:
			arg.WriteRune(c)
			inArg = true
		}
	}

	if quote != 0 {
		return nil, fmt.Errorf("unterminated %c quote", quote)
	}
	if escaped {
		return nil, errors.New("nothing follows the trailing backslash")
	}
	if inArg {
		args = append(args, arg.String())
	}
	return args, nil
}
