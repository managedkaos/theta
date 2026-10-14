package viewer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderFileDefaultHeadAndTail(t *testing.T) {
	path := writeTempFile(t, 12)

	out, err := RenderFile(path, Options{Head: 10, Tail: 10, Theme: "plain"})
	if err != nil {
		t.Fatalf("RenderFile returned error: %v", err)
	}

	if !strings.HasPrefix(out, path+"\n") {
		t.Fatalf("output should start with filename, got:\n%s", out)
	}
	if !strings.Contains(out, "\n"+separator+"\n") {
		t.Fatalf("output should contain separator, got:\n%s", out)
	}
	if !strings.Contains(out, "     1 | line 01\n") {
		t.Fatalf("output should include numbered head lines, got:\n%s", out)
	}
	if !strings.Contains(out, "    12 | line 12\n") {
		t.Fatalf("output should include numbered tail lines, got:\n%s", out)
	}
}

func TestRenderFileCustomHeadAndTail(t *testing.T) {
	path := writeTempFile(t, 8)

	out, err := RenderFile(path, Options{Head: 2, Tail: 3, Theme: "plain"})
	if err != nil {
		t.Fatalf("RenderFile returned error: %v", err)
	}

	want := path + "\n" +
		"     1 | line 01\n" +
		"     2 | line 02\n" +
		separator + "\n" +
		"     6 | line 06\n" +
		"     7 | line 07\n" +
		"     8 | line 08\n"

	if out != want {
		t.Fatalf("unexpected output\nwant:\n%s\ngot:\n%s", want, out)
	}
}

func TestRenderFileWithColor(t *testing.T) {
	path := writeTempFile(t, 1)

	out, err := RenderFile(path, Options{Head: 1, Tail: 1, Theme: "default", UseColor: true})
	if err != nil {
		t.Fatalf("RenderFile returned error: %v", err)
	}

	if !strings.Contains(out, "\033[1;36m"+path+"\033[0m") {
		t.Fatalf("filename should be colorized, got:\n%q", out)
	}
	if !strings.Contains(out, "\033[2;32m     1\033[0m | line 01") {
		t.Fatalf("line number should be colorized, got:\n%q", out)
	}
}

func writeTempFile(t *testing.T, lines int) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "sample.txt")
	var b strings.Builder
	for i := 1; i <= lines; i++ {
		fmt.Fprintf(&b, "line %02d\n", i)
	}

	if err := os.WriteFile(path, []byte(b.String()), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}

func TestEmptyFile(t *testing.T) {
	path := writeTempFile(t, 0)
	got, err := RenderFile(path, Options{Head: 10, Tail: 10, Theme: "plain"})
	if err != nil || got != path+"\n"+separator+"\n" {
		t.Fatalf("got %q, err %v", got, err)
	}
}

func TestLongLinesAndUnterminatedFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "long.txt")
	content := strings.Repeat("x", 128*1024)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := RenderFile(path, Options{Head: 1, Tail: 0, Theme: "plain"})
	if err != nil || got != path+"\n     1 | "+content+"\n"+separator+"\n" {
		t.Fatalf("unexpected long line output; err %v", err)
	}
}

func TestReadFailure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "large.txt")
	if err := os.WriteFile(path, []byte(strings.Repeat("x", 2*1024*1024)), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := RenderFile(path, Options{Head: 1, Tail: 1, Theme: "plain"})
	if err == nil || got != "" {
		t.Fatalf("expected scanner failure without partial output, got err %v", err)
	}
}
