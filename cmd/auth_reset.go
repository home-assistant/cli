package cmd

import (
	"fmt"

	helper "github.com/home-assistant/cli/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type User struct {
	Username  string
	Name      string
	Owner     bool
	Active    bool
	LocalOnly bool
}

var MsgLimitedAccess = "For security reasons, this command works only from the operating system terminal."

func getUsers() ([]User, error) {
	resp, err := helper.GenericJSONGet("auth", "list")

	if err != nil {
		return nil, err
	}

	var data *helper.Response
	var result []User

	if resp.IsSuccess() {
		var data *helper.Response = resp.Result().(*helper.Response)

		if data.Result != "ok" {
			err := fmt.Errorf("error returned from Supervisor: %s", data.Message)
			return nil, err
		}

		for _, user := range data.Data["users"].([]interface{}) {
			user := user.(map[string]interface{})
			result = append(result, User{
				Username:  user["username"].(string),
				Name:      user["name"].(string),
				Owner:     user["is_owner"].(bool),
				Active:    user["is_active"].(bool),
				LocalOnly: user["local_only"].(bool),
			})
		}
	} else {
		data = resp.Error().(*helper.Response)
		err := fmt.Errorf("error returned from Supervisor: %s", data.Message)
		return nil, err
	}

	return result, nil
}

func listUsers(users []User) {
	for i, user := range users {
		fmt.Printf("- %d: %s (%s)\n", i+1, user.Username, user.Name)
		fmt.Printf("     owner: %t, active: %t, local only: %t\n", user.Owner, user.Active, user.LocalOnly)
	}
}

var authResetCmd = &cobra.Command{
	Use:     "reset",
	Aliases: []string{"rst", "change"},
	Short:   "Reset the password of a Home Assistant user.",
	Long: `
This command allows you to change a password of a Home Assistant user.
Please note, this command is limited due to security reasons, and will
only work on some locations. For example, the Operating System CLI.
`,
	Example: `
  ha authentication reset --username "JohnDoe" --password "123SuperSecret!"
  ha authentication reset --interactive
`,
	ValidArgsFunction: cobra.NoFileCompletions,
	Args:              cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithField("args", args).Debug("auth reset")

		section := "auth"
		command := "reset"

		options := make(map[string]any)

		for _, value := range []string{
			"username",
			"password",
		} {
			val, err := cmd.Flags().GetString(value)
			if val != "" && err == nil && cmd.Flags().Changed(value) {
				options[value] = val
			}
		}

		interactive, _ := cmd.Flags().GetBool("interactive")
		if interactive {
			// prompt for user if not given
			if options["username"] != nil && options["username"] != "" {
				fmt.Printf("Changing password for user '%s'\n", options["username"])
			} else {
				users, err := getUsers()
				if err != nil {
					cmd.PrintErrln(MsgLimitedAccess)
					fmt.Println(err)
					ExitWithError = true
					return
				}

				fmt.Println("List of users:")
				listUsers(users)
				read, idx := helper.ReadInteger("Select a user to reset the password for", 3, 1, len(users))
				if read {
					user := users[idx-1]
					fmt.Printf("Changing password for user %d: %s (%s)\n", idx, user.Username, user.Name)
					options["username"] = user.Username
				} else {
					fmt.Println("Aborted.")
					ExitWithError = true
					return
				}
			}

			// prompt for password
			password, err := helper.ReadPassword(true)
			if err != nil {
				fmt.Printf("Failed to set password: %v\n", err)
				ExitWithError = true
				return
			}
			options["password"] = password
		}

		// change the password
		resp, err := helper.GenericJSONPost(section, command, options)
		if err != nil {
			cmd.PrintErrln(MsgLimitedAccess)
			fmt.Println(err)
			ExitWithError = true
		} else {
			ExitWithError = !helper.ShowJSONResponse(resp)
		}
	},
}

func init() {
	authResetCmd.Flags().String("username", "", "Username to reset the password for")
	authResetCmd.Flags().String("password", "", "New password to assign. Ignored in interactive mode")
	authResetCmd.Flags().Bool("interactive", false, "Use interactive prompt for entering username and/or password")
	authResetCmd.MarkFlagsOneRequired("username", "interactive")
	authResetCmd.MarkFlagsOneRequired("password", "interactive")
	authResetCmd.RegisterFlagCompletionFunc("username", cobra.NoFileCompletions)
	authResetCmd.RegisterFlagCompletionFunc("password", cobra.NoFileCompletions)
	authCmd.AddCommand(authResetCmd)
}
