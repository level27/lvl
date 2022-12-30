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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/level27/l27-go"
	"github.com/level27/lvl/utils"
	"github.com/spf13/cobra"
)

//
// common.go:
// Contains common functionality used by many commands.
//

// Contains parameters passed to get commands, like filter and max number of entries.
var optGetParameters l27.CommonGetParams

// Add common flags for get-style commands, such as --number and --filter.
func addCommonGetFlags(cmd *cobra.Command) {
	pf := cmd.Flags()

	pf.Int32VarP(&optGetParameters.Limit, "number", "n", optGetParameters.Limit, "How many things should we retrieve from the API?")
	pf.StringVarP(&optGetParameters.Filter, "filter", "f", optGetParameters.Filter, "How to filter API results?")
}

// Common flag to skip deletion confirmation prompts. Add flag with addDeleteConfirmFlag
var optDeleteConfirmed bool

func addDeleteConfirmFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&optDeleteConfirmed, "yes", "y", false, "Confirm deletion of entity without prompt")
}

// Common flag to wait on an operation like create, delete, etc...
var optWait bool

func addWaitFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&optWait, "wait", false, "Wait for the operation to finish on the API.")
}

// check for valid ID as type INT.
func checkSingleIntID(arg string, entity string) (l27.IntID, error) {
	id, err := l27.ParseID(arg)
	if err != nil {
		return 0, fmt.Errorf("not a valid %v ID: '%s'", entity, arg)
	}
	return id, nil
}

// Try to split the given cmd args into ID's (works with whitespace and komma's)
func CheckForMultipleIDs(ids []string) []string {
	var currIDs []string

	for _, id := range ids {
		tempID := strings.Split(id, ",")
		currIDs = append(currIDs, tempID...)
	}

	return currIDs
}

// --------------------------- DYNAMICALY SETTING PARAMETERS
// this can be used when you need flags to set different paramaters but you don't
// know the amount of parameters beforehand.
// we can use something like a parameter flag where the user has to set the keys and values himself.
// example lvl system cookbooks create -p key=value -p key2=value2
// function used for commands with dynamic parameters. (different parameters defined by 1 flag)
func SplitCustomParameters(customP []string) (map[string]interface{}, error) {
	checkedParameters := make(map[string]interface{})
	// loop over raw data set by user with -p flag
	for _, setParameter := range customP {
		// check if correct way is used to define parameters -> key=value
		if !strings.Contains(setParameter, "=") {
			// when there is no '=' in the parameter -> error
			return nil, fmt.Errorf("wrong way of defining parameter is used for: '%v'. (use:[ -p key=value ])", setParameter)
		}

		// split each parameter set by user into its key and value. put them in the dictionary
		line := strings.Split(setParameter, "=")
		// some keys can use multiple values. check if values seperated by comma
		if strings.Contains(line[1], ",") {
			values := strings.Split(line[1], ",")
			//removing spaces from splitted values
			for i := range values {
				values[i] = strings.Trim(values[i], " ")
			}
			// add key value pair to dict
			checkedParameters[strings.Trim(line[0], " ")] = values
		} else {
			// add key value pair to dict
			checkedParameters[strings.Trim(line[0], " ")] = strings.Trim(line[1], " ")
		}
	}

	return checkedParameters, nil
}

// Open a file passed as an argument.
// This handles the convention of "-" opening stdin.
func openArgFile(file string) (io.ReadCloser, error) {
	if file == "-" {
		return os.Stdin, nil
	}

	return os.Open(file)
}

// Get the value of an argument. If the argument begins with "@",
// it is interpreted as a file name and the contents of said file are read instead.
func readArgFileSupported(arg string) (string, error) {
	if strings.HasPrefix(arg, "@") {
		filename := arg[1:]
		file, err := openArgFile(filename)
		if err != nil {
			return "", fmt.Errorf("error while opening file %s: %v", filename, err)
		}

		contents, err := io.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("error while reading file %s: %v", filename, err)
		}

		return string(contents), nil
	}

	return arg, nil
}

// Load JSON settings from arg-specified file and merge it with override settings from other args.
func loadMergeSettings(fileName string, override map[string]interface{}) (map[string]interface{}, error) {
	jsonSettings, err := loadSettings(fileName)
	if err != nil {
		return nil, err
	}

	return mergeMaps(jsonSettings, override), nil
}

func loadSettings(fileName string) (map[string]interface{}, error) {
	if fileName == "" {
		return make(map[string]interface{}), nil
	}

	file, err := openArgFile(fileName)
	if err != nil {
		return nil, err
	}

	defer func() { file.Close() }()

	jsonBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var jsonSettings map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonSettings)
	if err != nil {
		return nil, err
	}

	return jsonSettings, nil
}

func mergeMaps(base map[string]interface{}, override map[string]interface{}) map[string]interface{} {
	var newMap = map[string]interface{}{}

	for k, v := range base {
		newMap[k] = v
	}

	for k, v := range override {
		newMap[k] = v
	}

	return newMap
}

func mergeSettingsWithEntity(entity interface{}, settings map[string]interface{}) map[string]interface{} {
	data := utils.RoundTripJson(entity).(map[string]interface{})
	return mergeMaps(data, settings)
}

var updateSettings = map[string]interface{}{}
var updateSettingsFile string

func settingsFileFlag(c *cobra.Command) {
	c.Flags().StringVarP(&updateSettingsFile, "settings-file", "s", "", "JSON file to read settings from. Pass '-' to read from stdin.")
}

// Add a string setting flag to a command, that will be stored in a map.
// This is intended to be easily used with PATCH APIs.
func settingString(c *cobra.Command, settings map[string]interface{}, name string, usage string) {
	c.Flags().Var(&stringMapValue{Map: settings, Name: name}, name, usage)
}

// Add a string setting flag to a command, that will be stored in a map. Shorthand version.
// This is intended to be easily used with PATCH APIs.
func settingStringP(c *cobra.Command, settings map[string]interface{}, name string, short string, usage string) {
	c.Flags().VarP(&stringMapValue{Map: settings, Name: name}, name, short, usage)
}

// Add an int setting flag to a command, that will be stored in a map.
// This is intended to be easily used with PATCH APIs.
func settingID(c *cobra.Command, settings map[string]interface{}, name string, usage string) {
	c.Flags().Var(&idMapValue{Map: settings, Name: name}, name, usage)
}

// Add an int setting flag to a command, that will be stored in a map. Shorthand version.
// This is intended to be easily used with PATCH APIs.
func settingIDP(c *cobra.Command, settings map[string]interface{}, name string, short string, usage string) {
	c.Flags().VarP(&idMapValue{Map: settings, Name: name}, name, short, usage)
}

// Add an int setting flag to a command, that will be stored in a map.
// This is intended to be easily used with PATCH APIs.
func settingInt32(c *cobra.Command, settings map[string]interface{}, name string, usage string) {
	c.Flags().Var(&int32MapValue{Map: settings, Name: name}, name, usage)
}

// Add an int setting flag to a command, that will be stored in a map. Shorthand version.
// This is intended to be easily used with PATCH APIs.
func settingInt32P(c *cobra.Command, settings map[string]interface{}, name string, short string, usage string) {
	c.Flags().VarP(&int32MapValue{Map: settings, Name: name}, name, short, usage)
}

// Add a bool setting flag to a command, that will be stored in a map.
// This is intended to be easily used with PATCH APIs.
func settingBool(c *cobra.Command, settings map[string]interface{}, name string, usage string) {
	c.Flags().Var(&boolMapValue{Map: settings, Name: name}, name, usage)
}

// Add a bool setting flag to a command, that will be stored in a map. Shorthand version.
// This is intended to be easily used with PATCH APIs.
func settingBoolP(c *cobra.Command, settings map[string]interface{}, name string, short string, usage string) {
	c.Flags().VarP(&boolMapValue{Map: settings, Name: name}, name, short, usage)
}

// Types below implement the pflag.Value interface for string/int/bool so that they can receive values assigned via command line.

type stringMapValue struct {
	Map  map[string]interface{}
	Name string
}

func (c *stringMapValue) String() string {
	val := c.Map[c.Name]
	if val == nil {
		return ""
	}

	return val.(string)
}

func (c *stringMapValue) Set(val string) error {
	c.Map[c.Name] = val
	return nil
}

func (c *stringMapValue) Type() string {
	return "string"
}

type idMapValue struct {
	Map  map[string]interface{}
	Name string
}

func (c *idMapValue) String() string {
	val := c.Map[c.Name]
	if val == nil {
		return ""
	}

	return fmt.Sprint(val)
}

func (c *idMapValue) Set(val string) error {
	i, err := l27.ParseID(val)
	if err != nil {
		return err
	}

	c.Map[c.Name] = i
	return nil
}

func (c *idMapValue) Type() string {
	return "ID"
}

type int32MapValue struct {
	Map  map[string]interface{}
	Name string
}

func (c *int32MapValue) String() string {
	val := c.Map[c.Name]
	if val == nil {
		return ""
	}

	return fmt.Sprint(val)
}

func (c *int32MapValue) Set(val string) error {
	i, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return err
	}

	c.Map[c.Name] = i
	return nil
}

func (c *int32MapValue) Type() string {
	return "int32"
}

type boolMapValue struct {
	Map  map[string]interface{}
	Name string
}

func (c *boolMapValue) String() string {
	val := c.Map[c.Name]
	if val == nil {
		return ""
	}

	return strconv.FormatBool(val.(bool))
}

func (c *boolMapValue) Set(val string) error {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}

	c.Map[c.Name] = b
	return nil
}

func (c *boolMapValue) Type() string {
	return "bool"
}

// Ask the user to confirm an action with a y/n prompt.
func confirmPrompt(message string) bool {
	for {
		fmt.Printf("%s [y]es/[n]o: ", message)
		var resp string
		_, err := fmt.Scan(&resp)
		if err != nil {
			log.Fatal(err)
		}
		resp = strings.ToLower(resp)

		if resp == "y" || resp == "yes" {
			return true
		} else if resp == "n" || resp == "no" {
			fmt.Printf("Operation cancelled\n")
			return false
		} else {
			continue
		}
	}
}

// Generic helper function to implement "get" commands.
// Takes in the list of arguments, a lookup function (name -> T*),
// a single-ID get function (int -> T) and a list function (get params -> []T)
// Will output the list of entities.
func resolveGets[T interface{}](
	args []string,
	lookup func(string) ([]T, error),
	getSingle func(l27.IntID) (T, error),
	getList func(l27.CommonGetParams) ([]T, error)) ([]T, error) {
	if len(args) == 0 {
		// No arguments, return full list from API.
		return getList(optGetParameters)
	} else {
		results := make([]T, 0, len(args))
		for _, val := range args {
			id, err := l27.ParseID(val)
			if err == nil {
				// Integer ID
				res, err := getSingle(id)
				if err != nil {
					return nil, err
				}

				results = append(results, res)
			} else {
				// Look up by name
				lookedUp, err := lookup(val)
				if err != nil {
					return nil, err
				}

				if len(lookedUp) == 0 {
					return nil, fmt.Errorf("unable to find '%s'", val)
				}

				results = append(results, lookedUp...)
			}
		}

		return results, nil
	}
}

func resolveShared[T interface{}](
	options []T,
	arg string,
	name string,
	getDesc func(T) string,
) (*T, error) {
	switch len(options) {
	case 0:
		return nil, fmt.Errorf("unable to find %s: %s", name, arg)
	case 1:
		return &options[0], nil
	default:
		// Multiple candidates, allow user to select which

		fmt.Printf("Multiple options exist for %s '%s':\n", name, arg)

		if !isStdinTerminal() {
			// If stdin isn't a terminal (e.g. being piped into) then we can't just prompt for input.
			// So abort in that case.
			return nil, errors.New("aborting because command not interactive")
		}

		for i, option := range options {
			fmt.Printf("[%d] %s\n", i, getDesc(option))
		}

		fmt.Printf("Choose one: ")
		var resp int
		_, err := fmt.Scan(&resp)
		if err != nil {
			return nil, err
		}

		if resp < 0 || resp >= len(options) {
			return nil, errors.New("invalid index given")
		}

		return &options[resp], nil
	}

}

func isStdinTerminal() bool {
	// See https://stackoverflow.com/a/43947435/4678631
	fi, _ := os.Stdin.Stat()

	return fi.Mode()&os.ModeCharDevice != 0
}

const waitPollInterval = 1 * time.Second
const waitPollTotal = 120

// Helper function to wait on the status of an entity to change to a desired value.
// poll is a function to fetch an entity from the API.
// status is a function to read the status field on the pulled entity.
// want is the desired entity state, and ignore is a set of "in-progress" states.
// For example, you may want to wait on a status of "ok" with an ignore of "updating".
// If the status were to change to "update_failed", the function returns with an error.
func waitForStatus[Entity any](
	poll func() (Entity, error),
	status func(Entity) string,
	want string,
	ignore []string,
) (Entity, error) {
	// poll and status are separate, to allow error handling to be done all in this function.

	var ent Entity
	var err error

	for i := 0; i < waitPollTotal; i += 1 {
		ent, err = poll()
		if err != nil {
			return ent, err
		}

		status := status(ent)

		if status == want {
			return ent, nil
		}

		if !sliceContains(ignore, status) {
			return ent, fmt.Errorf("got unexpected status: %s", status)
		}

		time.Sleep(waitPollInterval)
	}

	return ent, fmt.Errorf("timed out")
}

// Version of waitForStatus that waits on a system deletion.
// This means it waits for a status of "deleted" or a 404 error.
func waitForDelete[Entity any](
	poll func() (Entity, error),
	status func(Entity) string,
	ignore []string,
) error {
	// poll and status are separate, to allow error handling to be done all in this function.

	for i := 0; i < waitPollTotal; i += 1 {
		ent, err := poll()
		if err != nil {
			if errResp, ok := err.(l27.ErrorResponse); ok {
				if errResp.HTTPCode == 404 || errResp.Code == 404 {
					// 404 means deleted, we're done here.
					return nil
				}
			}
			return err
		}

		status := status(ent)

		if status == "deleted" {
			return nil
		}

		if !sliceContains(ignore, status) {
			return fmt.Errorf("got unexpected status: %s", status)
		}

		time.Sleep(waitPollInterval)
	}

	return fmt.Errorf("timed out")
}

func waitIndicator(block func()) {
	interval := 1 * time.Second

	done := make(chan bool)

	go func() {
		block()
		done <- true
	}()

	for {
		select {
		case <-done:
			fmt.Fprint(os.Stderr, "\n")
			return
		default:
			time.Sleep(interval)
			fmt.Fprint(os.Stderr, ".")
			continue
		}
	}
}
