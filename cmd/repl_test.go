package cmd

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestSplitLine(t *testing.T) {
	tests := []struct {
		line    string
		want    []string
		wantErr bool
	}{
		{line: "", want: nil},
		{line: "   ", want: nil},
		{line: "core info", want: []string{"core", "info"}},
		{line: "  core \t info  ", want: []string{"core", "info"}},
		{line: `backups new --name "My backup"`, want: []string{"backups", "new", "--name", "My backup"}},
		{line: `backups new --name 'a "b" $c \d'`, want: []string{"backups", "new", "--name", `a "b" $c \d`}},
		{line: `a\ b "c \" d" e\\`, want: []string{"a b", `c " d`, `e\`}},
		{line: `--password=pa'ss wo'rd`, want: []string{"--password=pass word"}},
		{line: `--name ""`, want: []string{"--name", ""}},
		{line: `"open`, wantErr: true},
		{line: `'open`, wantErr: true},
		{line: `trailing\`, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			got, err := splitLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseCompletions(t *testing.T) {
	tests := []struct {
		name        string
		output      string
		want        [][2]string
		wantNoSpace bool
	}{
		{
			name:   "candidates with and without description",
			output: "info\tShows information\nlogs\n_activeHelp_ a hint\n:4\n",
			want:   [][2]string{{"info", "Shows information"}, {"logs", ""}},
		},
		{name: "no space", output: "--name=\n:6\n", want: [][2]string{{"--name=", ""}}, wantNoSpace: true},
		{name: "error directive", output: "info\n:1\n"},
		{name: "file extensions are not candidates", output: "yaml\nyml\n:8\n"},
		{name: "directories are not candidates", output: "themes\n:16\n"},
		{name: "no directive", output: "Error: something failed\n"},
		{name: "empty", output: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, noSpace := parseCompletions(tt.output)
			if !reflect.DeepEqual(got, tt.want) || noSpace != tt.wantNoSpace {
				t.Errorf("got %q, %v, want %q, %v", got, noSpace, tt.want, tt.wantNoSpace)
			}
		})
	}
}

func TestCompleteLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		pos      int
		output   string
		wantArgs []string
		wantLine string
		wantPos  int
		wantList string
	}{
		{
			name:     "several matches are listed",
			line:     "core re",
			pos:      7,
			output:   "rebuild\tRebuild\nrestart\tRestart\ninfo\tInfo\n:4\n",
			wantArgs: []string{"core", "re"},
			wantLine: "core re",
			wantPos:  7,
			wantList: "rebuild  Rebuild\nrestart  Restart\n",
		},
		{
			name:     "common prefix is filled in",
			line:     "core r",
			pos:      6,
			output:   "rebuild\tRebuild\nrestart\tRestart\n:4\n",
			wantArgs: []string{"core", "r"},
			wantLine: "core re",
			wantPos:  7,
		},
		{
			name:     "unique match",
			line:     "core rest",
			pos:      9,
			output:   "rebuild\tRebuild\nrestart\tRestart\n:4\n",
			wantArgs: []string{"core", "rest"},
			wantLine: "core restart ",
			wantPos:  13,
		},
		{
			name:     "cursor inside the line",
			line:     "core rest --force",
			pos:      9,
			output:   "restart\tRestart\n:4\n",
			wantArgs: []string{"core", "rest"},
			wantLine: "core restart  --force",
			wantPos:  13,
		},
		{
			name:     "new word",
			line:     "core ",
			pos:      5,
			output:   "info\tInfo\n:4\n",
			wantArgs: []string{"core", ""},
			wantLine: "core info ",
			wantPos:  10,
		},
		{
			name:     "leading ha is dropped",
			line:     "ha core in",
			pos:      10,
			output:   "info\tInfo\n:4\n",
			wantArgs: []string{"core", "in"},
			wantLine: "ha core info ",
			wantPos:  13,
		},
		{
			name:     "flag value",
			line:     "--log-level=D",
			pos:      13,
			output:   "DEBUG\nINFO\n:4\n",
			wantArgs: []string{"--log-level=D"},
			wantLine: "--log-level=DEBUG ",
			wantPos:  18,
		},
		{
			name:     "no space directive",
			line:     "core restart --saf",
			pos:      18,
			output:   "--safe-mode=\tSafe mode\n:6\n",
			wantArgs: []string{"core", "restart", "--saf"},
			wantLine: "core restart --safe-mode=",
			wantPos:  25,
		},
		{
			name:     "builtins are offered for the first word",
			line:     "lo",
			pos:      2,
			output:   "logs\tLogs\n:4\n",
			wantArgs: []string{"lo"},
			wantLine: "log",
			wantPos:  3,
		},
		{
			name:     "builtins are not offered for later words",
			line:     "core lo",
			pos:      7,
			output:   "logs\tLogs\n:4\n",
			wantArgs: []string{"core", "lo"},
			wantLine: "core logs ",
			wantPos:  10,
		},
		{
			name:     "no match",
			line:     "core x",
			pos:      6,
			output:   "info\tInfo\n:4\n",
			wantArgs: []string{"core", "x"},
			wantLine: "core x",
			wantPos:  6,
		},
		{
			name:     "quoted lines are left alone",
			line:     `backups new --name "a `,
			pos:      22,
			wantLine: `backups new --name "a `,
			wantPos:  22,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotArgs []string
			r := &repl{complete: func(_ context.Context, args []string) string {
				gotArgs = args
				return tt.output
			}}
			line, pos, list := r.completeLine(t.Context(), tt.line, tt.pos)
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("completed %q, want %q", gotArgs, tt.wantArgs)
			}
			if line != tt.wantLine || pos != tt.wantPos {
				t.Errorf("got %q at %d, want %q at %d", line, pos, tt.wantLine, tt.wantPos)
			}
			if list != tt.wantList {
				t.Errorf("got list %q, want %q", list, tt.wantList)
			}
		})
	}
}

func TestReplLoop(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		banner   bool
		wantCode int
		wantRuns [][]string
		wantOut  string
	}{
		{
			name:     "commands run until exit",
			input:    "core info\r\rha host info\rexit\rcore stop\r",
			wantRuns: [][]string{{"core", "info"}, {"host", "info"}},
			wantOut:  "Leading 'ha' is not necessary",
		},
		{
			name:     "login leaves with its own exit code",
			input:    "login\r",
			wantCode: replLoginExitCode,
		},
		{
			name:     "quoted arguments",
			input:    "backups new --name 'My backup'\rexit\r",
			wantRuns: [][]string{{"backups", "new", "--name", "My backup"}},
		},
		{
			name:     "a line which cannot be split is not run",
			input:    "core 'info\rexit\r",
			wantRuns: nil,
		},
		{
			name:     "help mentions login",
			input:    "help core\rhelp\rexit\r",
			wantRuns: [][]string{{"help", "core"}, {"help"}},
			wantOut:  "Use \"login\" to enter operating system shell",
		},
		{
			name:     "Ctrl+C drops the line being edited",
			input:    "core stop\x03core info\rexit\r",
			wantRuns: [][]string{{"core", "info"}},
			wantOut:  "^C",
		},
		{
			name:     "Ctrl+D resets the prompt",
			input:    "core info\r\x04exit\r",
			banner:   true,
			wantRuns: [][]string{{"banner"}, {"core", "info"}, {"banner"}},
			wantOut:  replClearScreen,
		},
		{
			name:     "Ctrl+D inside a line does not reset",
			input:    "core info\x04\rexit\r",
			wantRuns: [][]string{{"core", "info"}},
		},
		{
			name:     "end of input leaves the prompt",
			input:    "core info\r",
			banner:   true,
			wantRuns: [][]string{{"banner"}, {"core", "info"}},
		},
		{
			name:     "Tab completes the line",
			input:    "co\tin\t\rexit\r",
			wantRuns: [][]string{{"core", "info"}},
		},
		{
			name:     "Ctrl+C ends a completion which hangs",
			input:    "hang\t\x03core info\rexit\r",
			wantRuns: [][]string{{"core", "info"}},
			wantOut:  "^C",
		},
		{
			name:     "keys typed during a completion are kept",
			input:    "slow\tcore info\rexit\r",
			wantRuns: [][]string{{"slow", "core", "info"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			var runs [][]string
			in := strings.NewReader(tt.input)
			r := &repl{
				in:     in,
				out:    &out,
				banner: tt.banner,
				run:    func(args []string) { runs = append(runs, args) },
				complete: func(ctx context.Context, args []string) string {
					switch {
					case args[0] == "hang":
						<-ctx.Done()
						return ""
					case args[0] == "slow":
						// Leave time to read the keys behind the Tab
						time.Sleep(20 * time.Millisecond)
						return "slow\tSlow\n:4\n"
					case len(args) == 1:
						return "core\tCore\n:4\n"
					}
					return "info\tInfo\n:4\n"
				},
			}
			if code := r.loop(); code != tt.wantCode {
				t.Errorf("got exit code %d, want %d", code, tt.wantCode)
			}
			if !reflect.DeepEqual(runs, tt.wantRuns) {
				t.Errorf("ran %q, want %q", runs, tt.wantRuns)
			}
			if !strings.Contains(out.String(), tt.wantOut) {
				t.Errorf("output %q does not contain %q", out.String(), tt.wantOut)
			}
		})
	}
}

func TestReplHistory(t *testing.T) {
	h := &replHistory{}
	for _, line := range []string{"core info", "", "  ", "host info", "host info", "core info"} {
		h.Add(line)
	}
	got := make([]string, 0, h.Len())
	for i := range h.Len() {
		got = append(got, h.At(i))
	}
	want := []string{"core info", "host info", "core info"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestReplLeavesInputToCommand(t *testing.T) {
	// The input is pasted at once, the second line is meant for the command
	in := strings.NewReader("auth reset\rsecret\rexit\r")
	var out bytes.Buffer
	var runs [][]string
	r := &repl{
		in:  in,
		out: &out,
		run: func(args []string) {
			runs = append(runs, args)
			line := make([]byte, len("secret\r"))
			if _, err := io.ReadFull(in, line); err != nil || string(line) != "secret\r" {
				t.Errorf("command read %q, %v", line, err)
			}
		},
	}
	if code := r.loop(); code != 0 {
		t.Errorf("got exit code %d, want 0", code)
	}
	if want := [][]string{{"auth", "reset"}}; !reflect.DeepEqual(runs, want) {
		t.Errorf("ran %q, want %q", runs, want)
	}
	if strings.Contains(out.String(), "secret") {
		t.Errorf("output %q holds the input of the command", out.String())
	}
}
