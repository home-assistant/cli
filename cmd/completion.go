package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:    "completion",
	Short:  "Generates bash/zsh completion scripts",
	Hidden: true,
	Long: `To load completion run

For Bash: . <(ha completion)
For ZSH: . <(ha completion --zsh)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc , ~/.profile
. <(ha completion)

# ~/.zshrc
. <(ha completion --zsh)
`,
	Run: func(cmd *cobra.Command, args []string) {

		_, err := cmd.Flags().GetBool("zsh")
		if err == nil && cmd.Flags().Changed("zsh") {
			rootCmd.GenZshCompletion(os.Stdout)

			// Workaround for the missing compdef in spf13's Cobra
			// ZSH completion generation.
			os.Stdout.WriteString(
				fmt.Sprintf(
					"\ncompdef _%s %s\n",
					path.Base(os.Args[0]),
					path.Base(os.Args[0]),
				),
			)
		} else {
			rootCmd.GenBashCompletionV2(os.Stdout, true)
		}
	},
}

func init() {
	completionCmd.Flags().Bool("zsh", false, "Generate ZSH completion script")
	rootCmd.AddCommand(completionCmd)
}
