package cmd

import (
	"fmt"
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
	cmd.Flags().String("device", "", "Device to mount, e.g. /dev/sda1 (disk mounts only)")
	cmd.Flags().String("uuid", "", "Filesystem UUID of the device to mount (disk mounts only)")
	cmd.Flags().Bool("read-only", false, "Is mount read-only (not available for backup mounts)")

	cmd.Flags().Lookup("read-only").NoOptDefVal = "true"
	cmd.MarkFlagsMutuallyExclusive("device", "uuid")
	cmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"cifs", "nfs", "disk"}, cobra.ShellCompDirectiveNoFileComp
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
	cmd.RegisterFlagCompletionFunc("device", mountsDeviceCompletions)
	cmd.RegisterFlagCompletionFunc("uuid", cobra.NoFileCompletions)
	cmd.RegisterFlagCompletionFunc("read-only", boolCompletions)
}

// Copy only flags the user set. Sending an untouched default would rewrite
// type or usage on update.
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
		"device",
		"uuid",
	} {
		val, err := cmd.Flags().GetString(value)
		if val != "" && err == nil && cmd.Flags().Changed(value) {
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

// Like mountFlagsToOptions, but include type/usage defaults for a new mount.
func mountFlagsToNewOptions(cmd *cobra.Command, options map[string]any) {
	mountFlagsToOptions(cmd, options)

	for _, value := range []string{"type", "usage"} {
		if _, ok := options[value]; ok {
			continue
		}
		if val, err := cmd.Flags().GetString(value); err == nil && val != "" {
			options[value] = val
		}
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
				if s, ok = m["uuid"].(string); ok && s != "" {
					ds = append(ds, s)
				}
				if s, ok = m["filesystem"].(string); ok && s != "" {
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

// Completions for --device from the host disks. Empty or 404 yields none.
func mountsDeviceCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := helper.GenericJSONGet("host", "disks")
	if err != nil || !resp.IsSuccess() {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ret []string
	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["disks"] != nil {
		if disks, ok := data.Data["disks"].([]any); ok {
			for _, disk := range disks {
				var d map[string]any
				if d, ok = disk.(map[string]any); !ok {
					continue
				}
				partitions, _ := d["partitions"].([]any)
				for _, partition := range partitions {
					var p map[string]any
					if p, ok = partition.(map[string]any); !ok {
						continue
					}
					if mountable, _ := p["mountable"].(bool); !mountable {
						continue
					}
					var device string
					if device, ok = p["device"].(string); !ok || device == "" {
						continue
					}
					ret = append(ret, device)
					if ds := mountDeviceDescription(d, p); ds != "" {
						ret[len(ret)-1] += "\t" + ds
					}
				}
			}
		}
	}
	return ret, cobra.ShellCompDirectiveNoFileComp
}

func mountDeviceDescription(disk, partition map[string]any) string {
	var ds []string
	if s, ok := disk["name"].(string); ok && s != "" {
		ds = append(ds, s)
	} else {
		var name []string
		for _, key := range []string{"vendor", "model"} {
			if s, ok := disk[key].(string); ok && s != "" {
				name = append(name, s)
			}
		}
		if len(name) != 0 {
			ds = append(ds, strings.Join(name, " "))
		}
	}
	if s, ok := partition["label"].(string); ok && s != "" {
		ds = append(ds, s)
	}
	if size, ok := partition["size"].(float64); ok && size > 0 {
		ds = append(ds, humanizeMountSize(size))
	}
	return strings.Join(ds, ", ")
}

// Decimal units, matching how drive capacity is labelled.
func humanizeMountSize(size float64) string {
	units := []string{"B", "kB", "MB", "GB", "TB"}
	i := 0
	for size >= 1000 && i < len(units)-1 {
		size /= 1000
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%.0f %s", size, units[i])
	}
	return fmt.Sprintf("%.1f %s", size, units[i])
}
