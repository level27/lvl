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
	"github.com/spf13/cobra"
)

var optNumber = "20"
var optFilter = ""

func addCommonGetFlags(cmd *cobra.Command) {
	pf := cmd.Flags()

	pf.StringVarP(&optNumber, "number", "n", optNumber, "How many things should we retrieve from the API?")
	pf.StringVarP(&optFilter, "filter", "f", optFilter, "How to filter API results?")
}
