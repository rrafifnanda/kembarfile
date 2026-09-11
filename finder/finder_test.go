package finder

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFind(t *testing.T) {
	root := t.TempDir()
	paths := []string{
		filepath.Join(root, "a.txt"),
		filepath.Join(root, "b.txt"),
		filepath.Join(root, "nested", "d.txt"),
	}
	if err := os.Mkdir(filepath.Join(root, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		if err := os.WriteFile(path, []byte("same"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "different.txt"), []byte("nope"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := Find(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Groups) != 1 {
		t.Fatalf("got %d groups, want 1", len(result.Groups))
	}
	if !reflect.DeepEqual(result.Groups[0].Paths, paths) {
		t.Fatalf("got paths %v, want %v", result.Groups[0].Paths, paths)
	}
	if result.Groups[0].Size != 4 || result.Reclaimable != 8 {
		t.Fatalf("got size=%d reclaimable=%d, want 4 and 8", result.Groups[0].Size, result.Reclaimable)
	}
}
