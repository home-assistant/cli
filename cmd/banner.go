package cmd

import (
	"fmt"
	"strings"
	"time"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const haBanner = `
       ▄██▄           _   _                                    
     ▄██████▄        | | | | ___  _ __ ___   ___               
   ▄████▀▀████▄      | |_| |/ _ \| '_ ` + "`" + ` _ \ / _ \              
 ▄█████    █████▄    |  _  | (_) | | | | | |  __/              
▄██████▄  ▄██████▄   |_| |_|\___/|_| |_| |_|\___|          _   
████████  ██▀  ▀██      / \   ___ ___(_)___| |_ __ _ _ __ | |_ 
███▀▀███  ██   ▄██     / _ \ / __/ __| / __| __/ _` + "`" + ` | '_ \| __|
██    ██  ▀ ▄█████    / ___ \\__ \__ \ \__ \ || (_| | | | | |_ 
███▄▄ ▀█  ▄███████   /_/   \_\___/___/_|___/\__\__,_|_| |_|\__|
▀█████▄   ███████▀

Welcome to the Home Assistant command line interface.
`

func supervisorGet(section string, command string) (outdata *(map[string]any), err error) {
	resp, err := helper.GenericJSONGet(section, command)
	if err != nil {
		return nil, err
	}

	var data *helper.Response
	if resp.IsSuccess() {
		data = resp.Result().(*helper.Response)
	} else {
		data = resp.Error().(*helper.Response)
	}
	if data.Result == "ok" {
		if len(data.Data) > 0 {
			outdata = &(data.Data)
		}
	} else {
		return nil, fmt.Errorf("error returned from Supervisor: %s", data.Message)
	}
	return outdata, nil
}

func getAddresses(addresses []any) string {
	addresses_str := make([]string, len(addresses))
	for i, v := range addresses {
		addresses_str[i] = fmt.Sprint(v)
	}
	return strings.Join(addresses_str, ", ")
}

var bannerCmd = &cobra.Command{
	Use:     "banner",
	Aliases: []string{"ba"},
	Short:   "Prints the CLI Home Assistant banner along with some useful information",
	Example: `
  ha banner
	`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("info")

		fmt.Print(haBanner)
		fmt.Println()

		nowait, err := cmd.Flags().GetBool("no-wait")
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		}

		var netinfo *(map[string]any)
		if !nowait {
			fmt.Println("Waiting for Supervisor to start. Press any key to skip waiting...")

			keyboardInterrupt, cancelKeyboard := helper.WaitForKeyboardInterrupt()
			defer cancelKeyboard()                 // Ensure terminal is always restored
			timeout := time.After(3 * time.Minute) // 3 minutes timeout
			tick := time.Tick(1 * time.Second)

			// Try immediately first, then wait for ticks
			firstAttempt := make(chan time.Time, 1)
			firstAttempt <- time.Now()

			checkSupervisor := func() bool {
				// We could use ping here, but Supervisor loads networking data asynchronously.
				// Networking info are very useful, so wait until networking data have been
				// loaded...
				var err error
				netinfo, err = supervisorGet("network", "info")
				if err == nil && netinfo != nil {
					netifaces, exist := (*netinfo)["interfaces"]
					if exist && len(netifaces.([]any)) > 0 {
						fmt.Println("Home Assistant Supervisor is running!")
						cancelKeyboard() // Restore terminal before continuing
						return true
					}
				}
				return false
			}

		waitLoop:
			for {
				select {
				case <-keyboardInterrupt:
					fmt.Println("Waiting interrupted by user.")
					return
				case <-timeout:
					fmt.Println("Supervisor is taking longer than expected to start. Use 'supervisor logs' to check logs.")
					return
				case <-firstAttempt:
					if checkSupervisor() {
						break waitLoop
					}
				case <-tick:
					if checkSupervisor() {
						break waitLoop
					}
				}
			}
		}

		fmt.Println("System information:")
		// If we don't have netinfo from the wait loop (nowait mode), try to fetch it now
		if netinfo == nil {
			var err error
			netinfo, err = supervisorGet("network", "info")
			if err != nil {
				fmt.Printf("  Network information unavailable: %s\n", err)
				return
			}
		}

		// Print network address information
		if netinfo != nil {
			netifaces, exist := (*netinfo)["interfaces"]
			if exist {
				for _, netiface := range netifaces.([]any) {
					nf := netiface.(map[string]any)
					title_ipv4 := fmt.Sprintf("IPv4 addresses for %s:", nf["interface"])
					title_ipv6 := fmt.Sprintf("IPv6 addresses for %s:", nf["interface"])

					if nf["ipv4"] == nil {
						fmt.Printf("  %-25s (No address)\n", title_ipv4)
					} else {
						ipv4 := nf["ipv4"].(map[string]any)
						ipv4_addresses := ipv4["address"].([]any)
						if len(ipv4_addresses) > 0 {
							fmt.Printf("  %-25s %s\n", title_ipv4, getAddresses(ipv4_addresses))
						} else {
							fmt.Printf("  %-25s (No address)\n", title_ipv4)
						}
					}

					if nf["ipv6"] != nil {
						ipv6 := nf["ipv6"].(map[string]any)
						ipv6_addresses := ipv6["address"].([]any)
						if len(ipv6_addresses) > 0 {
							fmt.Printf("  %-25s %s\n", title_ipv6, getAddresses(ipv6_addresses))
						}
					}
				}
			} else {
				fmt.Printf("  (No networking information)\n")
			}
		} else {
			fmt.Printf("  (Network information currently unavailable)\n")
		}
		fmt.Println()

		// Print Host URL
		hostinfo, err := supervisorGet("host", "info")
		if err != nil {
			ExitWithError = true
			fmt.Printf("  Host information unavailable: %s\n", err)
			return
		}
		if hostinfo == nil {
			return
		}

		coreinfo, err := supervisorGet("core", "info")
		if err != nil {
			ExitWithError = true
			fmt.Printf("  Core information unavailable: %s\n", err)
			return
		}
		if coreinfo == nil {
			return
		}

		protocol := "http"
		if ssl, ok := (*coreinfo)["ssl"].(string); ok && ssl == "true" {
			protocol = "https"
		}

		if os_version, ok := (*hostinfo)["operating_system"].(string); ok {
			fmt.Printf("  %-25s %s\n", "OS Version:", os_version)
		}
		if core_version, ok := (*coreinfo)["version"].(string); ok {
			fmt.Printf("  %-25s %s\n", "Home Assistant Core:", core_version)
		}
		fmt.Println()

		hostname, hostname_ok := (*hostinfo)["hostname"].(string)
		port, port_ok := (*coreinfo)["port"].(float64)

		if hostname_ok && port_ok {
			fmt.Printf("  %-25s %s://%s.local:%d\n", "Home Assistant URL:", protocol, hostname, int(port))
			fmt.Printf("  %-25s http://%s.local:%d\n", "Observer URL:", hostname, 4357)
		}
	},
}

func init() {
	rootCmd.AddCommand(bannerCmd)
	bannerCmd.Flags().Bool("no-wait", false, "Don't wait until Supervisor is started")
}
