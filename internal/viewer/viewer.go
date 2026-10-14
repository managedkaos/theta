// Package viewer renders numbered head and tail sections of text files.
package viewer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Options controls the line counts and styling of a rendered file.
type Options struct {
	Head, Tail int
	Theme      string
	UseColor   bool
}

const separator = "--------------------"

type theme struct {
	reset  string
	file   string
	lineNo string
	sep    string
}

var themes = map[string]theme{
	"default": {
		reset:  "\033[0m",
		file:   "\033[1;36m",
		lineNo: "\033[2;32m",
		sep:    "\033[2;33m",
	},
	"plain": {},
}

// RenderFile returns a complete file section, or an error without partial output.
func RenderFile(path string, opts Options) (string, error) {
	lines, err := readLines(path)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	t := themes[opts.Theme]

	writeColor(&b, opts, t.file, path)
	b.WriteByte('\n')

	headEnd := min(opts.Head, len(lines))
	for i := 0; i < headEnd; i++ {
		writeLine(&b, opts, i+1, lines[i])
	}

	writeColor(&b, opts, t.sep, separator)
	b.WriteByte('\n')

	tailStart := max(len(lines)-opts.Tail, 0)
	for i := tailStart; i < len(lines); i++ {
		writeLine(&b, opts, i+1, lines[i])
	}

	return b.String(), nil
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func writeLine(b *strings.Builder, opts Options, number int, text string) {
	t := themes[opts.Theme]
	if opts.UseColor {
		fmt.Fprintf(b, "%s%6d%s | %s\n", t.lineNo, number, t.reset, text)
		return
	}
	fmt.Fprintf(b, "%6d | %s\n", number, text)
}

func writeColor(b *strings.Builder, opts Options, code string, text string) {
	t := themes[opts.Theme]
	if opts.UseColor && code != "" {
		b.WriteString(code)
		b.WriteString(text)
		b.WriteString(t.reset)
		return
	}
	b.WriteString(text)
}
