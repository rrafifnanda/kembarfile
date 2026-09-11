# KembarFile MVP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a read-only Go CLI that finds duplicate files recursively and reports reclaimable bytes.

**Architecture:** A small `finder` package walks files, groups candidates by size, and confirms matches with SHA-256. The CLI only validates arguments and formats the reusable package result, leaving the core ready for a later desktop frontend.

**Tech Stack:** Go 1.27 standard library only

**Spec:** `docs/superpowers/specs/2026-09-11-kembarfile-design.md`

## Global Constraints

- The MVP is read-only: never delete, move, modify, or follow symbolic links.
- Use no external dependencies.
- Hash only size groups containing at least two regular files.
- Continue after unreadable nested entries, but fail when the root is invalid.
- Produce deterministic path ordering.

---

### Task 1: Complete KembarFile CLI

**Files:**
- Create: `go.mod`
- Create: `finder/finder_test.go`
- Create: `finder/finder.go`
- Create: `cmd/kembarfile/main.go`

**Interfaces:**
- Consumes: one root-directory path from the CLI.
- Produces: `finder.Find(root string) (finder.Result, error)` where `Result` contains `Groups []Group`, `Warnings []Warning`, and `Reclaimable int64`.

- [ ] **Step 1: Declare the module**

Create `go.mod`:

```go
module github.com/rrafifnanda/kembarfile

go 1.27
```

- [ ] **Step 2: Write the failing finder test**

Create `finder/finder_test.go`:

```go
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
```

- [ ] **Step 3: Run the test and confirm the expected failure**

Run: `go test ./finder`

Expected: compilation fails because `Find` is undefined.

- [ ] **Step 4: Implement the minimum reusable finder**

Create `finder/finder.go`:

```go
package finder

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

type Group struct {
	Size  int64
	Paths []string
}

type Warning struct {
	Path string
	Err  error
}

type Result struct {
	Groups      []Group
	Warnings    []Warning
	Reclaimable int64
}

func Find(root string) (Result, error) {
	info, err := os.Stat(root)
	if err != nil {
		return Result{}, err
	}
	if !info.IsDir() {
		return Result{}, fmt.Errorf("%s is not a directory", root)
	}

	var result Result
	bySize := make(map[int64][]string)
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if path == root {
				return walkErr
			}
			result.Warnings = append(result.Warnings, Warning{Path: path, Err: walkErr})
			return nil
		}
		if entry.Type()&fs.ModeSymlink != 0 {
			return nil
		}
		entryInfo, infoErr := entry.Info()
		if infoErr != nil {
			result.Warnings = append(result.Warnings, Warning{Path: path, Err: infoErr})
			return nil
		}
		if entryInfo.Mode().IsRegular() {
			bySize[entryInfo.Size()] = append(bySize[entryInfo.Size()], path)
		}
		return nil
	})
	if err != nil {
		return Result{}, err
	}

	for size, paths := range bySize {
		if len(paths) < 2 {
			continue
		}
		byHash := make(map[[sha256.Size]byte][]string)
		for _, path := range paths {
			file, openErr := os.Open(path)
			if openErr != nil {
				result.Warnings = append(result.Warnings, Warning{Path: path, Err: openErr})
				continue
			}
			hash := sha256.New()
			_, copyErr := io.Copy(hash, file)
			closeErr := file.Close()
			if copyErr != nil {
				result.Warnings = append(result.Warnings, Warning{Path: path, Err: copyErr})
				continue
			}
			if closeErr != nil {
				result.Warnings = append(result.Warnings, Warning{Path: path, Err: closeErr})
				continue
			}
			var digest [sha256.Size]byte
			copy(digest[:], hash.Sum(nil))
			byHash[digest] = append(byHash[digest], path)
		}
		for _, matches := range byHash {
			if len(matches) < 2 {
				continue
			}
			sort.Strings(matches)
			result.Groups = append(result.Groups, Group{Size: size, Paths: matches})
			result.Reclaimable += size * int64(len(matches)-1)
		}
	}
	sort.Slice(result.Groups, func(i, j int) bool {
		return result.Groups[i].Paths[0] < result.Groups[j].Paths[0]
	})
	return result, nil
}
```

- [ ] **Step 5: Format and run the finder test**

Run: `gofmt -w finder/finder.go finder/finder_test.go`

Run: `go test ./finder`

Expected: `ok github.com/rrafifnanda/kembarfile/finder`.

- [ ] **Step 6: Implement the thin CLI**

Create `cmd/kembarfile/main.go`:

```go
package main

import (
	"fmt"
	"os"

	"github.com/rrafifnanda/kembarfile/finder"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <folder>\n", os.Args[0])
		os.Exit(2)
	}

	result, err := finder.Find(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	for _, warning := range result.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s: %v\n", warning.Path, warning.Err)
	}
	for i, group := range result.Groups {
		fmt.Printf("Group %d (%d bytes):\n", i+1, group.Size)
		for _, path := range group.Paths {
			fmt.Println(" ", path)
		}
	}
	fmt.Printf("%d duplicate group(s), %d bytes reclaimable\n", len(result.Groups), result.Reclaimable)
}
```

- [ ] **Step 7: Verify the complete program**

Run: `gofmt -w cmd/kembarfile/main.go`

Run: `go test ./...`

Expected: all packages pass.

Run: `go vet ./...`

Expected: no output and exit status 0.

Create a temporary sample with shell commands, then run `go run ./cmd/kembarfile <sample-folder>`.

Expected: one duplicate group, all matching paths, and the correct reclaimable-byte total.

- [ ] **Step 8: Commit the working MVP**

```bash
git add go.mod finder/finder.go finder/finder_test.go cmd/kembarfile/main.go
git commit -m "feat: add duplicate file finder CLI"
```
