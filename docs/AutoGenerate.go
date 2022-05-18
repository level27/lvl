package docs

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func GenerateDocumentation(cmd *cobra.Command, dir string, filePrepender, linkHandler func(string) string) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := GenerateDocumentation(c, dir, filePrepender, linkHandler); err != nil {
			return err
		}
	}

	category := strings.Split(cmd.CommandPath(), " ")
	if len(category) > 1{
		dir = ResolveDirectory(category[1])
	}

	basename := strings.Replace(cmd.CommandPath(), " ", "_", -1) + ".md"

	// split the above filename and check the second word
	// the second word is the name of the command, this way we can automate the files
	// to go in right folder for each sort of command.
	
	
	filename := filepath.Join(dir, basename)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.WriteString(f, filePrepender(filename)); err != nil {
		return err
	}
	if err := doc.GenMarkdownCustom(cmd, f, linkHandler); err != nil {
		return err
	}
	return nil
}

func ResolveDirectory(command string) string{
	var dir string
	switch command {
	case "app":
		dir = "docs/Apps"
	case "domain":
		dir = "docs/Domains"
	case "component":
		dir = "docs/Components"
	case "job":
		dir = "docs/Jobs"
	case "mail":
		dir = "docs/Mails"
	case "network":
		dir = "docs/Networks"
	case "system":
		dir = "docs/Systems"
	case "systemgroup":
		dir = "docs/Systemgroups"
	case "organisation":
		dir = "docs/Organisations"
	case "login":
		dir = "docs/Login"
	case "region":
		dir = "docs/Regions"
	default:
		dir = "docs/"
	}
	return dir
}

