package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash/zsh completion scripts",
	Hidden: true,
	Long: `To load completion run

For Bash: . <(hassio-cli completion)
For ZSH: . <(hassio-cli completion --zsh)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc , ~/.profile
. <(hassio-cli completion)

# ~/.zshrc
. <(hassio-cli completion --zsh)
`,
	Run: func(cmd *cobra.Command, args []string) {

		_, err := cmd.Flags().GetBool("zsh")
		if err == nil && cmd.Flags().Changed("zsh") {
			rootCmd.GenZshCompletion(os.Stdout)
		} else {
			rootCmd.GenBashCompletion(os.Stdout)
		}
	},
}

func init() {
	completionCmd.Flags().Bool("zsh", false, "Generate ZSH completion script")
	rootCmd.AddCommand(completionCmd)
}
