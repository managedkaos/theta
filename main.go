package main

import (
	"fmt"
	"github.com/managedkaos/theta/cmd"
	"os"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "theta:", err)
		os.Exit(1)
	}
}
