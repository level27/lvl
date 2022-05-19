package main

import (
	"log"

	"bitbucket.org/level27/lvl/cmd"
)



	func main(){

	// this will generate de documentation and put everything the correct folder. 
	// this updated Docs folder is the folder that needs to replace the docs folder at cli.docs.level27.eu
	
		cmd.RootCmd.DisableAutoGenTag =true
		err := GenerateDocumentation(cmd.RootCmd, "docs", func(s string) string { return "" }, func(s string) string { return s })
		if err != nil {
		log.Fatal(err)
	}
	}