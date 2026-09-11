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
