package cmd

import (
	"strings"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var mountsCmd = &cobra.Command{
	Use:     "mounts",
	Aliases: []string{"mount", "mnts", "mnt"},
	Short:   "Get information, update or configure mounts in Supervisor",
	Long: `
The mounts command allows you to manage mounts in Supervisor by exposing
commands to view, mount, update or remove mounts such as network shares.`,
	Example: `
  ha mounts info
  ha mounts add my_share --usage media --type cifs --server server.local --share media`,
}

func init() {
	rootCmd.AddCommand(mountsCmd)
}

func addMountFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("type", "t", "cifs", "Type of mount")
	cmd.Flags().StringP("usage", "u", "media", "Usage of mount within Home Assistant")
	cmd.Flags().StringP("server", "s", "", "IP address or hostname of network share server")
	cmd.Flags().IntP("port", "o", 0, "Port to use if network share is exposed on non-standard port for the type")
	cmd.Flags().StringP("share", "r", "", "Share to mount (cifs mounts only)")
	cmd.Flags().StringP("username", "n", "", "Username to use for authentication (cifs mounts only)")
	cmd.Flags().StringP("password", "p", "", "Password to use for authentication (cifs mounts only)")
	cmd.Flags().StringP("version", "v", "", "Version to use for the mount (cifs mounts only)")
	cmd.Flags().StringP("path", "a", "", "Path to mount (nfs mounts only)")
	cmd.Flags().Bool("read-only", false, "Is mount read-only (not available for backup mounts)")

	cmd.Flags().Lookup("read-only").NoOptDefVal = "true"
	cmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"cifs", "nfs"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc("usage", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"backup", "media", "share"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.RegisterFlagCompletionFunc("server", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("port", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("share", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("username", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("password", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("version", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("path", cobra.NoFileCompletions)
}

func mountFlagsToOptions(cmd *cobra.Command, options map[string]any) {
	for _, value := range []string{
		"type",
		"usage",
		"server",
		"share",
		"path",
		"username",
		"password",
		"version",
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

	roVal, roErr := cmd.Flags().GetBool("read-only")
	if roErr == nil && cmd.Flags().Changed("read-only") {
		options["read_only"] = roVal
	}
}

func mountsCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	resp, err := helper.GenericJSONGet("mounts", "")
	if err != nil || !resp.IsSuccess() {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ret []string
	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["mounts"] != nil {
		if mounts, ok := data.Data["mounts"].([]any); ok {
			for _, mount := range mounts {
				var m map[string]any
				if m, ok = mount.(map[string]any); !ok {
					continue
				}
				var s string
				if s, ok = m["name"].(string); !ok {
					continue
				}
				ret = append(ret, s)
				var ds []string
				if s, ok = m["state"].(string); ok && s != "" {
					ds = append(ds, s)
				}
				if s, ok = m["usage"].(string); ok && s != "" {
					ds = append(ds, s)
				}
				if s, ok = m["server"].(string); ok && s != "" {
					ds = append(ds, s)
				}
				if s, ok = m["share"].(string); ok && s != "" {
					ds = append(ds, s)
				}
				if s, ok = m["path"].(string); ok && s != "" {
					ds = append(ds, s)
				}
				if len(ds) != 0 {
					ret[len(ret)-1] += "\t" + strings.Join(ds, ", ")
				}
			}
		}
	}
	return ret, cobra.ShellCompDirectiveNoFileComp
}
