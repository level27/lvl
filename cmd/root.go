/*
Copyright © 2021 Level27 info@level27.be

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
	"fmt"
	"log"
	"os"

	"bitbucket.org/level27/lvl/utils"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var apiKey string
var level27Client *utils.Client

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lvl",
	Short: "CLI tool to manage Level27 entities",
	Long:  `lvl is a CLI tool that empowers users.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lvl.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "apikey", "", "API key")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("apikey", rootCmd.PersistentFlags().Lookup("apikey"))
	viper.BindPFlag("toggle", rootCmd.PersistentFlags().Lookup("toggle"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".lvl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".lvl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		apiKey = viper.GetString("apiKey")
		level27Client = utils.NewAPIClient("https://api.level27.eu/v1", apiKey)
	} else {
		// config file is not found we create it
		fmt.Println(cfgFile)
		if cfgFile == "" {
			home, _ := homedir.Dir()
			file := home + "/.lvl.yaml"
			_, err = os.Stat(file)
			if os.IsNotExist(err) {
				file, err := os.Create(file)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()
			}
		}
	}
}
