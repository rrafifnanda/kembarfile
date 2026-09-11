package main

import (
	"fmt"
	"io"
	"os"

	"github.com/rrafifnanda/kembarfile/finder"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: kembarfile <folder>")
		return 2
	}

	result, err := finder.Find(args[0])
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	for _, warning := range result.Warnings {
		fmt.Fprintf(stderr, "warning: %s: %v\n", warning.Path, warning.Err)
	}
	for i, group := range result.Groups {
		fmt.Fprintf(stdout, "Group %d (%d bytes):\n", i+1, group.Size)
		for _, path := range group.Paths {
			fmt.Fprintln(stdout, " ", path)
		}
	}
	fmt.Fprintf(stdout, "%d duplicate group(s), %d bytes reclaimable\n", len(result.Groups), result.Reclaimable)
	return 0
}
