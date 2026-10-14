package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func execute(args ...string) (string, error) {
	if args == nil {
		args = []string{}
	}
	root := NewRootCommand()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs(args)
	err := root.Execute()
	return out.String(), err
}

func sample(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "sample.txt")
	if err := os.WriteFile(path, []byte("one\ntwo\nthree\nfour\n"), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestFlags(t *testing.T) {
	path := sample(t)
	want := path + "\n     1 | one\n--------------------\n     3 | three\n     4 | four\n"
	for _, args := range [][]string{
		{"--lines", "3", "--head", "1", "--tail", "2", "--color", "never", "--theme", "plain", path},
		{"-l", "3", "-h", "1", "-t", "2", "-c", "never", "-T", "plain", path},
		{path, "-h", "1", "-t", "2", "-l", "3"},
	} {
		got, err := execute(args...)
		if err != nil || got != want {
			t.Fatalf("%v: got %q, err %v; want %q", args, got, err, want)
		}
	}
}

func TestInvalidFlags(t *testing.T) {
	for _, args := range [][]string{
		{"--lines=-1"}, {"--head=-1"}, {"--tail=-1"},
		{"--color=invalid"}, {"--theme=invalid"}, {"--unknown"}, {"-h"}, {"-t", "no"},
	} {
		if _, err := execute(args...); err == nil {
			t.Errorf("%v should fail", args)
		}
	}
}

func TestHelpAndVersion(t *testing.T) {
	for _, args := range [][]string{nil, {"--help"}} {
		got, err := execute(args...)
		if err != nil || !strings.Contains(got, "theta [flags] FILE") {
			t.Fatalf("help: %q, %v", got, err)
		}
	}
	for _, flag := range []string{"--version", "-v"} {
		got, err := execute(flag)
		if err != nil || !strings.HasPrefix(got, "theta "+Version+"\n") {
			t.Fatalf("version: %q, %v", got, err)
		}
	}
}

func TestColors(t *testing.T) {
	path := sample(t)
	for _, tc := range []struct {
		mode, theme string
		color       bool
	}{
		{"auto", "default", false}, {"never", "default", false},
		{"always", "default", true}, {"always", "plain", false},
	} {
		got, err := execute("-c", tc.mode, "-T", tc.theme, path)
		if err != nil || strings.Contains(got, "\x1b[") != tc.color {
			t.Fatalf("%+v: %q, %v", tc, got, err)
		}
	}
}

func TestMultipleFilesAndFailure(t *testing.T) {
	path := sample(t)
	missing := filepath.Join(t.TempDir(), "missing")
	got, err := execute("-l", "0", path, missing, path)
	section := path + "\n--------------------\n"
	if got != section+"\n"+section || err == nil || !strings.Contains(err.Error(), missing) {
		t.Fatalf("got %q, err %v", got, err)
	}
}

func TestDefaultAndZeroCounts(t *testing.T) {
	path := sample(t)
	got, err := execute(path)
	if err != nil || strings.Count(got, "one\n") != 2 || strings.Count(got, "four\n") != 2 {
		t.Fatalf("default: %q, %v", got, err)
	}
	got, err = execute("-h", "0", "-t", "1", path)
	want := path + "\n--------------------\n     4 | four\n"
	if err != nil || got != want {
		t.Fatalf("zero: %q, %v", got, err)
	}
}

func TestDashFilename(t *testing.T) {
	got, err := execute("-l", "0", "--", "--help")
	if err == nil || got != "" || !strings.Contains(err.Error(), "--help:") {
		t.Fatalf("got %q, err %v", got, err)
	}
}
