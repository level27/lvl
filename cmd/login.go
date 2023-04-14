/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/common-nighthawk/go-figure"
	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// loginCmd represents the login command
var optLoginUsername string
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to CP4",
	Long: `Log in to CP4
You will be prompted for your user name (usually email, like example@level27.be) and password.
If --username is specified, user name will not be prompted, only password. This can be useful for scripts.

After the first login, further logins will suggest you to re-use the previous username when logging in.
`,
	Example: `Log in manually:
lvl login

Log in automatically:
cat password.txt | lvl login -u my.awesome.email@level27.be`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var login l27.Login
		var err error
		username, password, err := loginPromptCredentials()
		if err != nil {
			return err
		}

		client := makeApiClient(apiUrl, "")

		request := l27.LoginRequest{
			Username:   username,
			Password:   password,
			TwoFAToken: viper.GetString("2faKey"),
		}

		login, err = client.Login2FA(&request)
		if err != nil {
			// Check if it's a 2FA-related error we can handle appropriately.
			var l27err l27.ErrorResponse
			var ok bool
			if l27err, ok = err.(l27.ErrorResponse); !ok {
				return err
			}

			if l27err.Message == "2FA is requested" {
				// User pressed the button on the control panel
				// to try to set up 2FA but didn't finish yet.
				fmt.Printf("You are currently setting up 2-factor authorization. Please finish setting up 2FA on the website before trying to log in.\n")
				return nil
			}

			if l27err.Message == "2FA is required" {
				// 2FA required by organisation but not set up yet by user.
				fmt.Printf("Your organisation requires you to set up 2-factor authentication before logging in again. Please proceed on the website to do this.\n")
				return nil
			}

			if l27err.Message == "6digitCode is invalid" || l27err.Message == "2FA token is invalid" {
				// "2FA token is invalid": 2FA is enabled, "Trust this device" token is invalid
				// "6digitCode is invalid": 2FA is enabled
				// In both cases we need to re-enter a 2FA code from the authenticator.
				login, err = login2fa(client, username, password)
				if err != nil {
					return err
				}
			}
		}

		fmt.Println()
		loginFigure := figure.NewColorFigure("LEVEL27 CLI", "", "gray", true)
		loginFigure.Print()
		fmt.Println()
		fmt.Printf("Successfully logged in using: %s\n", username)

		utils.SaveConfig("apikey", login.Hash)
		utils.SaveConfig("2faKey", login.Hash2FA)
		utils.SaveConfig("user_id", login.User.ID)
		utils.SaveConfig("org_id", login.User.Organisation.ID)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&optLoginUsername, "username", "u", "", "User name to log in with")
}

func loginPromptCredentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)
	var err error

	var username string
	if optLoginUsername != "" {
		username = optLoginUsername
	} else {
		lastUsername := viper.GetString("last_username")
		prompt := "Enter Username"
		if lastUsername != "" {
			prompt += fmt.Sprintf(" (empty for %s)", lastUsername)
		}
		prompt += ": "

		fmt.Print(prompt)
		username, err = reader.ReadString('\n')
		username = strings.TrimSpace(username)
		if username == "" && lastUsername != "" {
			username = lastUsername
		}

		utils.SaveConfig("last_username", username)

		if err != nil {
			return "", "", err
		}
	}

	var bytePassword []byte
	if term.IsTerminal(int(syscall.Stdin)) {
		fmt.Print("Enter Password: ")
		bytePassword, err = term.ReadPassword(int(syscall.Stdin))
		// So that the next line of output doesn't overlap the password prompt's former spot
		fmt.Println()
	} else {
		bytePassword, err = reader.ReadBytes('\n')
	}

	if err != nil {
		return "", "", err
	}

	password := strings.TrimSpace(string(bytePassword))
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

// 1st return value: 2FA code
// 2nd return value: "trust this device?"
func loginPrompt2fa() (string, bool, error) {
	if !term.IsTerminal(int(syscall.Stdin)) {
		return "", false, errors.New("cannot prompt 2FA if piping input. Log in manually and trust device instead to avoid 2FA prompts and automate login")
	}

	reader := bufio.NewReader(os.Stdin)

	// Prompt 2FA code
	var code string
	var err error
	for {
		fmt.Print("Enter 2FA authenticator code: ")
		code, err = reader.ReadString('\n')
		if err != nil {
			return "", false, err
		}

		code = strings.TrimSpace(code)
		if len(code) != 6 {
			fmt.Println("Code must be exactly 6 digits")
			continue
		}

		break
	}

	fmt.Printf("Trust this device? [y]es/[n]o (default: no): ")
	resp, err := reader.ReadString('\n')
	if err != nil {
		return "", false, err
	}

	resp = strings.TrimSpace(resp)
	resp = strings.ToLower(resp)

	trustDevice := resp == "y" || resp == "yes"

	return code, trustDevice, nil
}

func login2fa(client *l27.Client, username string, password string) (l27.Login, error) {
	code, trust, err := loginPrompt2fa()
	if err != nil {
		return l27.Login{}, err
	}

	request := l27.LoginRequest{
		Username:        username,
		Password:        password,
		TrustThisDevice: trust,
		SixDigitCode:    code,
	}

	login, err := client.Login2FA(&request)
	return login, err
}
