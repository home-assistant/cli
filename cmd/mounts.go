package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var mountsCmd = &cobra.Command{
	Use:     "mounts",
	Aliases: []string{"mount", "mnts", "mnt"},
	Short:   "Get information, update or configure mounts in Supervisor",
	Long: `
The multicast command allows you to manage mounts in Supervisor by exposing
commands to view, mount, update or remove mounts such as network shares.`,
	Example: `
  ha mounts info
  ha mounts add my_share --usage media --type cifs --server server.local --share media`,
}

func init() {
	rootCmd.AddCommand(mountsCmd)
}

func portCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	vals := make([]string, 1, 65535)
	for i := 0; i <= 100; i++ {
		vals = append(vals, fmt.Sprint(i))
	}
	return vals, cobra.ShellCompDirectiveNoFileComp
}

func addMountFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("type", "t", "cifs", "Type of mount")
	cmd.Flags().StringP("usage", "u", "media", "Usage of mount within Home Assistant")
	cmd.Flags().StringP("server", "s", "", "IP address or hostname of network share server")
	cmd.Flags().IntP("port", "o", 0, "Port to use if network share is exposed on non-standard port for the type")
	cmd.Flags().StringP("share", "r", "", "Share to mount (cifs mounts only)")
	cmd.Flags().StringP("username", "u", "", "Username to use for authentication (cifs mounts only)")
	cmd.Flags().StringP("password", "p", "", "Password to use for authentication (cifs mounts only)")
	cmd.Flags().StringP("path", "h", "", "Path to mount (nfs mounts only)")

	cmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"cifs", "nfs"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc("usage", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"backup", "media"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc("server", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("port", portCompletions)
	cmd.RegisterFlagCompletionFunc("share", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("username", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("password", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("path", cobra.NoFileCompletions)
}

func mountFlagsToOptions(cmd *cobra.Command, options map[string]interface{}) {
	for _, value := range []string{
		"type",
		"usage",
		"server",
		"share",
		"path",
		"username",
		"password",
	} {
		val, err := cmd.Flags().GetString(value)
		if val != "" && err == nil {
			options[value] = val
		}
	}

	val, err := cmd.Flags().GetInt("port")
	if val > 0 && err == nil && cmd.Flags().Changed("port") {
		options["port"] = val
	}
}
