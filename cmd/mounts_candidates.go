package cmd

import (
	"fmt"
	"log/slog"
	"strings"

	helper "github.com/home-assistant/cli/client"
	"github.com/spf13/cobra"
)

var mountsCandidatesCmd = &cobra.Command{
	Use:     "candidates",
	Aliases: []string{"can", "cand"},
	Short:   "Shows local disks which could be mounted",
	Long: `
Shows the local block devices Supervisor considers usable as a disk mount,
along with details of the drive each one belongs to.
`,
	Example: `
  ha mounts candidates
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("mounts candidates", "args", args)

		section := "mounts"
		command := "candidates"

		resp, err := helper.GenericJSONGet(section, command)
		if err != nil {
			helper.PrintError(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

// mountsCandidatesDeviceCompletions offers the devices Supervisor reports as
// mountable, offering nothing at all if it cannot ask: a host without UDisks2
// returns an empty list, and a Supervisor predating the endpoint returns 404.
func mountsCandidatesDeviceCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	resp, err := helper.GenericJSONGet("mounts", "candidates")
	if err != nil || !resp.IsSuccess() {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var ret []string
	data := resp.Result().(*helper.Response)
	if data.Result == "ok" && data.Data["candidates"] != nil {
		if candidates, ok := data.Data["candidates"].([]any); ok {
			for _, candidate := range candidates {
				var c map[string]any
				if c, ok = candidate.(map[string]any); !ok {
					continue
				}
				var device string
				if device, ok = c["device"].(string); !ok || device == "" {
					continue
				}
				ret = append(ret, device)
				if ds := mountCandidateDescription(c); ds != "" {
					ret[len(ret)-1] += "\t" + ds
				}
			}
		}
	}
	return ret, cobra.ShellCompDirectiveNoFileComp
}

func mountCandidateDescription(candidate map[string]any) string {
	var ds []string
	if drive, ok := candidate["drive"].(map[string]any); ok {
		var name []string
		for _, key := range []string{"vendor", "model"} {
			if s, ok := drive[key].(string); ok && s != "" {
				name = append(name, s)
			}
		}
		if len(name) != 0 {
			ds = append(ds, strings.Join(name, " "))
		}
	}
	if s, ok := candidate["label"].(string); ok && s != "" {
		ds = append(ds, s)
	}
	if size, ok := candidate["size"].(float64); ok && size > 0 {
		ds = append(ds, humanizeMountSize(size))
	}
	return strings.Join(ds, ", ")
}

// humanizeMountSize uses decimal units, the way drive capacity is labelled.
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

func init() {
	mountsCmd.AddCommand(mountsCandidatesCmd)
}
