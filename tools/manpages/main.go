package main

import (
	"github.com/managedkaos/theta/cmd"
	"github.com/spf13/cobra/doc"
	"log"
	"os"
)

func main() {
	if err := os.MkdirAll("man", 0755); err != nil {
		log.Fatal(err)
	}
	root := cmd.NewRootCommand()
	root.DisableAutoGenTag = true
	if err := doc.GenManTree(root, &doc.GenManHeader{Title: "THETA", Section: "1"}, "man"); err != nil {
		log.Fatal(err)
	}
}
