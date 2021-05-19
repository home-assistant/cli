package cmd

import (
	"fmt"
	"strings"
	"time"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const haBanner = `
| |  | |                          /\           (_)   | |            | |  
| |__| | ___  _ __ ___   ___     /  \   ___ ___ _ ___| |_ __ _ _ __ | |_ 
|  __  |/ _ \| '_ \ _ \ / _ \   / /\ \ / __/ __| / __| __/ _\ | '_ \| __|
| |  | | (_) | | | | | |  __/  / ____ \\__ \__ \ \__ \ || (_| | | | | |_ 
|_|  |_|\___/|_| |_| |_|\___| /_/    \_\___/___/_|___/\__\__,_|_| |_|\__|

Welcome to the Home Assistant command line.
`

func supervisorGet(section string, command string) (outdata *(map[string]interface{}), err error) {
	base := viper.GetString("endpoint")

	resp, err := helper.GenericJSONGet(base, section, command)
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
		return nil, fmt.Errorf("Error returned from Supervisor: %s", data.Message)
	}
	return outdata, nil
}

func getAddresses(addresses []interface{}) string {
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
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("info")

		fmt.Print(haBanner)
		fmt.Println()

		nowait, err := cmd.Flags().GetBool("no-wait")
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
		}

		if !nowait {
			for i := 0; i < 60; i++ {
				// We could use ping here, but Supervisor loads networking data asynchronously.
				// Networking info are very useful, so wait until networking data have been
				// loaded...
				netinfo, err := supervisorGet("network", "info")
				if err == nil && netinfo != nil {

					netifaces, exist := (*netinfo)["interfaces"]
					if exist && len(netifaces.([]interface{})) > 0 {
						break
					}
				}
				if i == 0 {
					fmt.Println("Waiting for Supervisor to startup...")
				}
				time.Sleep(1 * time.Second)
			}
		}

		fmt.Println("System information")
		netinfo, err := supervisorGet("network", "info")
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}
		if netinfo == nil {
			return
		}

		// Print network address information
		netifaces, exist := (*netinfo)["interfaces"]
		if exist {
			for _, netiface := range netifaces.([]interface{}) {
				nf := netiface.(map[string]interface{})
				title_ipv4 := fmt.Sprintf("IPv4 addresses for %s:", nf["interface"])
				title_ipv6 := fmt.Sprintf("IPv6 addresses for %s:", nf["interface"])

				if nf["ipv4"] == nil {
					fmt.Printf("  %-25s (no address)\n", title_ipv4)
				} else {
					ipv4 := nf["ipv4"].(map[string]interface{})
					fmt.Printf("  %-25s %s\n", title_ipv4, getAddresses(ipv4["address"].([]interface{})))
				}

				if nf["ipv6"] != nil {
					ipv6 := nf["ipv6"].(map[string]interface{})
					ipv6_addresses := ipv6["address"].([]interface{})
					if len(ipv6_addresses) > 0 {
						fmt.Printf("  %-25s %s\n", title_ipv6, getAddresses(ipv6_addresses))
					}
				}
			}
		} else {
			fmt.Printf("  (no networking information)")
		}
		fmt.Println()

		// Print Host URL
		hostinfo, err := supervisorGet("host", "info")
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}
		if hostinfo == nil {
			return
		}
		coreinfo, err := supervisorGet("core", "info")
		if err != nil {
			fmt.Println(err)
			ExitWithError = true
			return
		}
		if coreinfo == nil {
			return
		}

		protocol := "http"
		if (*coreinfo)["ssl"] == "true" {
			protocol = "https"
		}

		port, _ := (*coreinfo)["port"].(float64)
		fmt.Printf("  %-25s %s\n", "OS Version:", (*hostinfo)["operating_system"])
		fmt.Printf("  %-25s %s\n", "Home Assistant Core:", (*coreinfo)["version"])
		fmt.Println()
		fmt.Printf("  %-25s %s://%s.local:%d\n", "Home Assistant URL:", protocol, (*hostinfo)["hostname"], int(port))
		fmt.Printf("  %-25s http://%s.local:%d\n", "Observer URL:", (*hostinfo)["hostname"], 4357)
	},
}

func init() {
	rootCmd.AddCommand(bannerCmd)
	bannerCmd.Flags().Bool("no-wait", false, "Don't wait until Supervisor is started")
}
