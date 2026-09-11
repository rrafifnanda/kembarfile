package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunReportsDuplicates(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("same"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var stdout, stderr bytes.Buffer

	if code := run([]string{root}, &stdout, &stderr); code != 0 {
		t.Fatalf("got exit code %d and stderr %q", code, stderr.String())
	}
	if output := stdout.String(); !strings.Contains(output, "1 duplicate group(s), 4 bytes reclaimable") {
		t.Fatalf("unexpected output %q", output)
	}
}
