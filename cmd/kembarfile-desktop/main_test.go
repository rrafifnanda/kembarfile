package main

import (
	"errors"
	"testing"

	"github.com/rrafifnanda/kembarfile/internal/finder"
)

func TestFormatResult(t *testing.T) {
	tests := []struct {
		name   string
		result finder.Result
		want   string
	}{
		{
			name: "duplicates and warning",
			result: finder.Result{
				Groups:      []finder.Group{{Size: 4, Paths: []string{"/a", "/b"}}},
				Warnings:    []finder.Warning{{Path: "/c", Err: errors.New("permission denied")}},
				Reclaimable: 4,
			},
			want: "Grup 1 (4 bytes)\n/a\n/b\n\nPeringatan\n/c: permission denied\n",
		},
		{
			name: "no duplicates",
			want: "Tidak ditemukan file duplikat.\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := formatResult(test.result); got != test.want {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}

func TestIconEmbedded(t *testing.T) {
	if len(icon) == 0 {
		t.Fatal("desktop icon is empty")
	}
}
