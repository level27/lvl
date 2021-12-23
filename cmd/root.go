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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"text/tabwriter"
	"text/template"

	"bitbucket.org/level27/lvl/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var apiKey string
var apiUrl string
var output string
var Level27Client *utils.Client

// NOTE: subcommands like get add themselves to root in their init().
// This requires importing them manually in main.go

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "lvl",
	Short: "CLI tool to manage Level27 entities",
	Long:  `lvl is a CLI tool that empowers users.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		outputSet := viper.GetString("output")
		if outputSet != "text" && outputSet != "json" && output != "yaml" {
			return fmt.Errorf("invalid output mode specified: '%s'", outputSet)
		}

		return nil
	},

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lvl.yaml)")
	RootCmd.PersistentFlags().StringVar(&apiKey, "apikey", "", "API key")
	RootCmd.PersistentFlags().BoolVar(&utils.TraceRequests, "trace", false, "Do detailed network request logging")
	RootCmd.PersistentFlags().StringVarP(&output, "output", "o", "text", "Specifies output mode for commands. Accepted values are text or json.")

	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("apikey", RootCmd.PersistentFlags().Lookup("apikey"))
	viper.BindPFlag("output", RootCmd.PersistentFlags().Lookup("output"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetDefault("apiUrl", "https://api.level27.eu/v1")

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
		apiUrl = viper.GetString("apiUrl")
		Level27Client = utils.NewAPIClient(apiUrl, apiKey)
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

// Output formatting functions

// Output tabular data from the CLI. Respects the --output flag.
// objects must be a slice of some set of objects.
// titles is the list of table headers,
// and fields contains the corresponding fields to read from the objects to fill said columns.
// When outputting as a structured format like JSON, the titles and fields are unused,
// and the slice is simply serialized directly.
func outputFormatTable(objects interface{}, titles []string, fields []string) {
	outputMode := viper.GetString("output")
	switch outputMode {
	case "text":
		outputFormatTableText(objects, titles, fields)
	case "json":
		outputFormatTableJson(objects)
	case "yaml":
		outputFormatTableYaml(objects)
	}
}

// Output templated data from the CLI (such as a describe output). Respects the --output flag.
// object must be the object to output
// templatePath must be the path to the go template formatting it under text mode.
// When outputting as a structured format like JSON,
// the template path is unused and the object is simply serialized directly.
func outputFormatTemplate(object interface{}, templatePath string) {
	outputMode := viper.GetString("output")
	switch outputMode {
	case "text":
		outputFormatTemplateText(object, templatePath)
	case "json":
		outputFormatTemplateJson(object)
	case "yaml":
		outputFormatTemplateYaml(object)
	}
}

func outputFormatTableText(objects interface{}, titles []string, fields []string) {
	// Have to use reflection for this because no generics in go (yet).
	s := reflect.ValueOf(objects)

	if s.Kind() != reflect.Slice {
		panic("outputFormatTable must be given a slice!")
	}

	w := tabwriter.NewWriter(os.Stdout, 4, 8, 4, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, strings.Join(titles, "\t"))

	for i := 0; i < s.Len(); i++ {
		val := s.Index(i)

		first := true
		for _, fieldName := range fields {
			fld := val.FieldByName(fieldName)

			if !first {
				fmt.Fprintf(w, "\t");
			}

			first = false;
			fmt.Fprintf(w, "%v", fld.Interface())
		}

		fmt.Fprintf(w, "\n")
	}
}

func outputFormatTableJson(objects interface{}) {
	b, _ := json.Marshal(objects)
	fmt.Println(string(b))
}

func outputFormatTableYaml(objects interface{}) {
	b, _ := yaml.Marshal(roundTripJson(objects))
	fmt.Println(string(b))
}

func outputFormatTemplateText(object interface{}, templatePath string) {
	_, fileName := path.Split(templatePath)
	tmpl := template.Must(template.New(fileName).Funcs(utils.TemplateHelpers()).ParseFiles(templatePath))
	err := tmpl.Execute(os.Stdout, object)
	if err != nil {
		panic(err)
	}
}

func outputFormatTemplateJson(object interface{}) {
	b, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func outputFormatTemplateYaml(object interface{}) {
	b, err := yaml.Marshal(roundTripJson(object))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func roundTripJson(obj interface{}) interface{} {
	// Round-trip through JSON so we use the JSON (camelcased) key names in the YAML without having to re-define them
	bJson, _ := json.Marshal(obj)
	var interf interface{}
	json.Unmarshal(bJson, &interf)
	return interf
}

func convertStringToId(id string) (int, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("'%s' is not a valid ID", id)
	}

	return intId, nil
}

func convertStringsToIds(ids []string) ([]int, error) {
	ints := make([]int, len(ids))
	for idx, id := range ids {
		intId, err := strconv.Atoi(id)
		if err != nil {
			return nil, fmt.Errorf("'%s' is not a valid ID", id)
		}
		ints[idx] = intId
	}

	return ints, nil
}