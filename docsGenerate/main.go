package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/level27/lvl/cmd"
)

func main() {
	// this will generate de documentation and put everything the correct folder.
	// this updated Docs folder is the folder that needs to replace the docs folder at cli.docs.level27.eu

	cmd.RootCmd.DisableAutoGenTag = true
	os.RemoveAll("docs/docsAuto")
	err := GenerateDocumentation(cmd.RootCmd, "docs/docsAuto", func(s string) string { return "" }, func(s string) string { return s })

	if err != nil {
		log.Fatal(err)
	}

	// Overlay docs/docs over docs/docsAuto so our manual docs go over.

	err = copyDir("docs/docs", "docs/docsAuto")
	if err != nil {
		log.Fatal(err)
	}
}

func copyDir(src string, dst string) error {
	files, err := os.ReadDir("docs/docs")
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Recursive so we can have more than home.md
	for _, file := range files {
		name := file.Name()
		err = copyFile(filepath.Join(src, name), filepath.Join(dst, name))
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
