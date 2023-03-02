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
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var apiKey string
var apiUrl string
var output string
var Level27Client *l27.Client

// NOTE: subcommands like get add themselves to root in their init().
// This requires importing them manually in main.go

var errSilent = errors.New("silentErr")

//go:embed version.txt
var version string

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:           "lvl",
	Short:         "CLI tool to manage Level27 entities",
	Long:          `lvl is a CLI tool that empowers users.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	Version:       strings.TrimSpace(version),

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		outputSet := viper.GetString("output")
		if outputSet != "text" && outputSet != "json" && output != "yaml" && output != "id" {
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
	// See https://github.com/spf13/cobra/issues/914 for some of the error handling details.

	err := RootCmd.Execute()
	if err != nil {
		if err != errSilent {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		}

		os.Exit(1)
	}
}

var traceRequests bool

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lvl.yaml)")
	RootCmd.PersistentFlags().StringVar(&apiKey, "apikey", "", "API key")
	RootCmd.PersistentFlags().BoolVar(&traceRequests, "trace", false, "Do detailed network request logging. This is intended for debugging and should not be parsed.")
	RootCmd.PersistentFlags().StringVarP(&output, "output", "o", "text", "Specifies output mode for commands. Accepted values are 'text', 'json', 'yaml' or 'id'.")

	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("apikey", RootCmd.PersistentFlags().Lookup("apikey"))
	viper.BindPFlag("output", RootCmd.PersistentFlags().Lookup("output"))

	RootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return errSilent
	})
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
	if err := viper.ReadInConfig(); err != nil {
		if !isViperConfigNotFound(err) {
			// Error we can't handle.
			log.Fatalf("Error while reading configuration: %s", err.Error())
		}

		// Config file doesn't exist yet, create it.
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

	// Load values from config.
	apiKey = viper.GetString("apiKey")
	apiUrl = viper.GetString("apiUrl")
	Level27Client = l27.NewAPIClient(apiUrl, apiKey)
	Level27Client.DefaultRequestHeaders["User-Agent"] = getUserAgent()
	if traceRequests {
		Level27Client.TraceRequests(&colorRequestTracer{})
	}
}

func getUserAgent() string {
	return fmt.Sprintf("level27_lvl/%s", strings.TrimSpace(version))
}

// Viper doesn't consistently return a ConfigFileNotFoundError in all cases.
// This helper function checks for ENOENT aswell.
func isViperConfigNotFound(err error) bool {
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		return true
	}

	if os.IsNotExist(err) {
		return true
	}

	return false
}

// Output formatting functions

// Output tabular data from the CLI. Respects the --output flag.
// objects must be a slice of some set of objects.
// titles is the list of table headers,
// and fields contains the corresponding fields names to read from the objects to fill said columns.
// Field names can contain "." separators to allow nested property field access.
// When outputting as a structured format like JSON, the titles and fields are unused,
// and the slice is simply serialized directly.
func outputFormatTable(objects interface{}, titles []string, fields []string) {
	outputMode := viper.GetString("output")
	switch outputMode {
	case "text":
		fieldsInterface := make([]interface{}, len(fields))
		for i := range fields {
			fieldsInterface[i] = fields[i]
		}
		outputFormatTableText(objects, titles, fieldsInterface)
	case "json":
		outputFormatTableJson(objects)
	case "yaml":
		outputFormatTableYaml(objects)
	case "id":
		outputFormatTableId(objects)
	}
}

// Equivalent to outputFormatTable, but takes in a slice of interfaces as field names instead.
// If a field is a string, it acts the same as outputFormatTable.
// If instead the field is a func with a single parameter and return value,
// it will be called with the row object to get the column value.
func outputFormatTableFuncs(objects interface{}, titles []string, fields []interface{}) {
	outputMode := viper.GetString("output")
	switch outputMode {
	case "text":
		outputFormatTableText(objects, titles, fields)
	case "json":
		outputFormatTableJson(objects)
	case "yaml":
		outputFormatTableYaml(objects)
	case "id":
		outputFormatTableId(objects)
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
	case "id":
		outputFormatTemplateId(object)
	}
}

func outputFormatTableText(objects interface{}, titles []string, fields []interface{}) {
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
		for _, field := range fields {
			var fld reflect.Value
			if fieldPath, isString := field.(string); isString {
				fld = val
				for _, fieldName := range strings.Split(fieldPath, ".") {
					fld = fld.FieldByName(fieldName)

				}
			} else {
				// Assume function that returns actal field value.
				fld = reflect.ValueOf(field).Call([]reflect.Value{val})[0]
			}

			if !first {
				fmt.Fprintf(w, "\t")
			}

			first = false
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
	b, _ := yaml.Marshal(utils.RoundTripJson(objects))
	fmt.Println(string(b))
}

func outputFormatTableId(objects interface{}) {
	s := reflect.ValueOf(objects)

	// Automatically get the ID field out of each struct in the list.

	if s.Kind() != reflect.Slice {
		panic("outputFormatTable must be given a slice!")
	}

	for i := 0; i < s.Len(); i++ {
		val := s.Index(i)
		if val.Kind() != reflect.Struct {
			break
		}

		fld := val.FieldByName("ID")
		if fld.IsZero() {
			break
		}

		fmt.Println(fld.Interface())
	}
}

//go:embed templates
var templates embed.FS

func outputFormatTemplateText(object interface{}, templatePath string) {
	_, fileName := path.Split(templatePath)

	tmpl := template.New(fileName)
	tmpl.Funcs(sprig.TxtFuncMap())
	tmpl.Funcs(utils.MakeTemplateHelpers(tmpl))
	tmpl = template.Must(tmpl.ParseFS(templates, templatePath))
	tmpl = template.Must(tmpl.ParseFS(templates, "templates/helpers/*.tmpl"))

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
	b, err := yaml.Marshal(utils.RoundTripJson(object))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func outputFormatTemplateId(object interface{}) {
	// Automatically get the ID field from the value.

	val := reflect.ValueOf(object)
	if val.Kind() != reflect.Struct {
		return
	}

	fld := val.FieldByName("ID")
	if fld.IsZero() {
		return
	}

	fmt.Println(fld.Interface())
}

// Tries to convert a string command line argument to an integer ID
func convertStringToID(id string) (l27.IntID, error) {
	intID, err := l27.ParseID(id)
	if err != nil {
		return 0, fmt.Errorf("'%s' is not a valid ID", id)
	}

	return intID, nil
}

// Tries to convert a slice of command line arguments to integer IDs
func convertStringsToIDs(ids []string) ([]l27.IntID, error) {
	ints := make([]l27.IntID, len(ids))
	for idx, id := range ids {
		intID, err := convertStringToID(id)
		if err != nil {
			return nil, err
		}
		ints[idx] = intID
	}

	return ints, nil
}

type colorRequestTracer struct{}

func (c *colorRequestTracer) TraceRequest(method string, url string, reqData []byte) {
	fmt.Fprintf(os.Stderr, "Request: %s %s\n", method, url)
	if len(reqData) != 0 {
		colored, err := utils.ColorJson(reqData)
		var str string
		if err == nil {
			str = string(colored)
		} else {
			str = string(reqData)
		}

		fmt.Fprintf(os.Stderr, "Request Body: %s\n", str)
	}
}

func (c *colorRequestTracer) TraceResponse(response *http.Response) {
	fmt.Fprintf(os.Stderr, "Response: %d %s\n", response.StatusCode, http.StatusText(response.StatusCode))
}

func (c *colorRequestTracer) TraceResponseBody(response *http.Response, data []byte) {
	bodyPrint := data
	if json.Valid(bodyPrint) {
		bodyPrint, _ = utils.ColorJson(bodyPrint)
	}
	fmt.Fprintf(os.Stderr, "Response Body: %s\n", string(bodyPrint))
}
