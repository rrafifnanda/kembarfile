package main

import (
	_ "embed"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/rrafifnanda/kembarfile/internal/finder"
)

//go:embed kembarfile.svg
var icon []byte

func main() {
	a := app.NewWithID("io.github.rrafifnanda.kembarfile")
	w := a.NewWindow("KembarFile")
	w.SetIcon(fyne.NewStaticResource("kembarfile.svg", icon))
	w.Resize(fyne.NewSize(900, 600))

	folderPath := widget.NewEntry()
	folderPath.SetPlaceHolder("Belum ada folder yang dipilih")
	folderPath.Disable()
	output := widget.NewMultiLineEntry()
	output.Wrapping = fyne.TextWrapOff
	output.SetText("Pilih folder untuk mulai mencari file duplikat.")
	output.Disable()
	status := widget.NewLabel("Siap")
	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	var selectedFolder string
	scanButton := widget.NewButton("Scan", nil)
	scanButton.Disable()
	chooseButton := widget.NewButton("Pilih Folder", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uri == nil {
				return
			}
			selectedFolder = uri.Path()
			folderPath.SetText(selectedFolder)
			scanButton.Enable()
			status.SetText("Folder siap dipindai")
		}, w).Show()
	})

	setBusy := func(busy bool) {
		if busy {
			chooseButton.Disable()
			scanButton.Disable()
			progress.Show()
			progress.Start()
			status.SetText("Memindai folder...")
			return
		}
		chooseButton.Enable()
		scanButton.Enable()
		progress.Stop()
		progress.Hide()
	}
	scanButton.OnTapped = func() {
		setBusy(true)
		root := selectedFolder
		go func() {
			result, err := finder.Find(root)
			fyne.Do(func() {
				setBusy(false)
				if err != nil {
					status.SetText("Pemindaian gagal")
					dialog.ShowError(err, w)
					return
				}
				output.SetText(formatResult(result))
				status.SetText(fmt.Sprintf("%d grup duplikat, %d bytes dapat dibebaskan", len(result.Groups), result.Reclaimable))
			})
		}()
	}

	header := container.NewVBox(
		widget.NewLabelWithStyle("KembarFile", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Cari file duplikat tanpa mengubah data Anda."),
		container.NewBorder(nil, nil, nil, chooseButton, folderPath),
		container.NewHBox(scanButton),
		progress,
	)
	w.SetContent(container.NewBorder(header, status, nil, nil, output))
	w.ShowAndRun()
}

func formatResult(result finder.Result) string {
	var text strings.Builder
	if len(result.Groups) == 0 {
		text.WriteString("Tidak ditemukan file duplikat.\n")
	}
	for i, group := range result.Groups {
		fmt.Fprintf(&text, "Grup %d (%d bytes)\n", i+1, group.Size)
		for _, path := range group.Paths {
			fmt.Fprintln(&text, path)
		}
		text.WriteByte('\n')
	}
	if len(result.Warnings) > 0 {
		text.WriteString("Peringatan\n")
		for _, warning := range result.Warnings {
			fmt.Fprintf(&text, "%s: %v\n", warning.Path, warning.Err)
		}
	}
	return text.String()
}
