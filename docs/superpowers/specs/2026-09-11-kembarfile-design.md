# KembarFile Design

## Goal

Build a safe, dependency-free Go CLI that recursively finds duplicate files and reports how much storage could be reclaimed. The first release never deletes or moves files. Its scanning core must remain reusable by a later desktop interface.

## CLI

The initial command is:

```text
kembarfile <folder>
```

The command validates that exactly one readable directory was supplied. It prints duplicate groups containing the file size and paths, followed by the number of duplicate groups and reclaimable bytes. Reclaimable space counts every copy after one retained copy in each group.

## Scanning and matching

The scanner walks the supplied directory recursively using the Go standard library. It considers regular files only and skips symbolic links, directories, sockets, devices, and other non-regular entries.

Files are grouped by byte size first. SHA-256 is calculated only for size groups containing at least two files. Files are duplicates only when both their size and full SHA-256 digest match. Results are sorted by path so repeated runs produce stable output.

An unreadable file does not abort the entire scan. The result includes its path and error so the CLI can print a warning to standard error and continue. Failure to access the root directory is fatal and returns a non-zero exit status.

## Structure

- `finder/finder.go`: reusable scan function and result types.
- `finder/finder_test.go`: focused duplicate-detection test using temporary files.
- `cmd/kembarfile/main.go`: argument validation, output formatting, and exit status.
- `go.mod`: module declaration; no external dependencies.

The finder package owns filesystem traversal and duplicate classification. The command owns presentation only. A future desktop program can import the finder package without running or parsing the CLI.

## Safety and limits

The MVP is read-only. It does not delete, move, modify, or follow symbolic links. Hard-linked paths may appear as duplicates because they have identical contents; inode-aware handling is deferred until real usage shows it is needed.

The scanner reads candidates sequentially. Parallel hashing, ignore patterns, progress displays, persistent indexes, and desktop UI are deferred until the CLI is correct and measured on real directories.

## Verification

One table-style test creates duplicate and distinct temporary files, runs the finder, and verifies the returned groups and reclaimable byte count. Final verification runs `go test ./...`, `go vet ./...`, and a manual CLI scan of a temporary sample directory.
