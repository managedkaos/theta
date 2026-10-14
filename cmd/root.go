// Package cmd defines theta's command-line interface.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/managedkaos/theta/internal/viewer"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// Build metadata is injected by Make or GoReleaser.
var (
	Version          = "dev"
	GitBranch        = "unknown"
	BuildDate        = "unknown"
	BuildEnvironment = "local"
)

// NewRootCommand constructs an independent command for each invocation.
func NewRootCommand() *cobra.Command {
	var lines, head, tail int
	var color, theme string
	var version bool
	root := &cobra.Command{
		Use:           "theta [flags] FILE [FILE...]",
		Short:         "Print the first and last lines of files with line numbers",
		Long:          "Theta prints the first and last lines of each file with line numbers, a 20-dash separator, and optional ANSI color.",
		Example:       "  theta notes.txt\n  theta -l 5 notes.txt report.txt\n  theta -h 20 -t 5 -c always app.log\n  theta --theme plain app.log",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(c *cobra.Command, files []string) error {
			if version {
				_, err := fmt.Fprintf(c.OutOrStdout(), "theta %s\nBranch: %s\nBuilt: %s\nEnvironment: %s\n", Version, GitBranch, BuildDate, BuildEnvironment)
				return err
			}
			if lines < 0 || (c.Flags().Changed("head") && head < 0) || (c.Flags().Changed("tail") && tail < 0) {
				return errors.New("--lines, --head, and --tail must be non-negative")
			}
			if color != "auto" && color != "always" && color != "never" {
				return fmt.Errorf("unknown color mode %q (use auto, always, or never)", color)
			}
			if theme != "default" && theme != "plain" {
				return fmt.Errorf("unknown theme %q (use default or plain)", theme)
			}
			if len(files) == 0 {
				return c.Help()
			}
			if !c.Flags().Changed("head") {
				head = lines
			}
			if !c.Flags().Changed("tail") {
				tail = lines
			}
			useColor := color == "always"
			if color == "auto" {
				if f, ok := c.OutOrStdout().(*os.File); ok {
					useColor = term.IsTerminal(int(f.Fd()))
				}
			}
			opts := viewer.Options{Head: head, Tail: tail, Theme: theme, UseColor: useColor && theme != "plain"}
			wrote := false
			var failures []error
			for _, file := range files {
				output, err := viewer.RenderFile(file, opts)
				if err != nil {
					failures = append(failures, fmt.Errorf("%s: %w", file, err))
					continue
				}
				if wrote {
					if _, err := fmt.Fprintln(c.OutOrStdout()); err != nil {
						return err
					}
				}
				if _, err := fmt.Fprint(c.OutOrStdout(), output); err != nil {
					return err
				}
				wrote = true
			}
			return errors.Join(failures...)
		},
	}
	root.CompletionOptions.DisableDefaultCmd = true
	f := root.Flags()
	f.IntVarP(&lines, "lines", "l", 10, "number of lines for both head and tail")
	f.IntVarP(&head, "head", "h", 10, "number of head lines, overriding --lines")
	f.IntVarP(&tail, "tail", "t", 10, "number of tail lines, overriding --lines")
	f.StringVarP(&color, "color", "c", "auto", "color output: auto, always, never")
	f.StringVarP(&theme, "theme", "T", "default", "theme: default, plain")
	f.Bool("help", false, "show help (--help; -h selects head lines)")
	f.BoolVarP(&version, "version", "v", false, "print version and build metadata")
	return root
}
