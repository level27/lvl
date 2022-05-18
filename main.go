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
package main

import (
	"log"

	"bitbucket.org/level27/lvl/cmd"
	"bitbucket.org/level27/lvl/docs"
	"github.com/spf13/cobra/doc"
)

func main() {
	// ---- // UNCOMMENT BELOW AND RUN MAIN.GO TO GENERATE MARKDOWN DOCS FOR WHOLE COMMAND TREE // ---- //
	cmd.RootCmd.DisableAutoGenTag =true
	err := docs.GenerateDocumentation(cmd.RootCmd, "docs", func(s string) string { return "" }, func(s string) string { return s })
	if err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}

func Test(){
	ok := doc.GenMarkdownTree(cmd.RootCmd, "docs")
log.Println(ok)
}
