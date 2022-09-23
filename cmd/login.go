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
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to CP4",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var login l27.Login
		username, password, _ := credentials()

		client := l27.NewAPIClient(apiUrl, "")
		login, err := client.Login(username, password)
		if err != nil {
			return err
		}

		fmt.Println()
		loginFigure := figure.NewColorFigure("LEVEL27 CLI", "", "gray", true)
		loginFigure.Print()
		fmt.Println()
		fmt.Printf("Successfully logged in using: %s\n", username)

		// fmt.Println(login.Hash)
		utils.SaveConfig("apikey", login.Hash)
		utils.SaveConfig("user_id", login.User.ID)
		utils.SaveConfig("org_id", login.User.Organisation.ID)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	lastUsername := viper.GetString("last_username")
	prompt := "Enter Username"
	if lastUsername != "" {
		prompt += fmt.Sprintf(" (empty for %s)", lastUsername)
	}
	prompt += ": "

	fmt.Print(prompt)
	username, err := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	if username == "" && lastUsername != "" {
		username = lastUsername
	}

	utils.SaveConfig("last_username", username)

	if err != nil {
		return "", "", err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	// So that the next line of output doesn't overlap the password prompt's former spot
	fmt.Println()
	if err != nil {
		return "", "", err
	}

	password := strings.TrimSpace(string(bytePassword))
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}
