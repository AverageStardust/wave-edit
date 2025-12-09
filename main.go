package main

import (
	"errors"
	"os"
	"wave-edit/riff"
	"wave-edit/wave"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var ErrExpectedWave = errors.New("expected WAVE file")

var mainWindow fyne.Window

func main() {
	app := app.New()
	mainWindow = app.NewWindow("WAVE edit")

	mainWindow.Resize(fyne.NewSize(640, 480))
	mainWindow.SetContent(widget.NewLabel("Drag and drop a WAVE file"))
	mainWindow.SetOnDropped(handleFile)
	mainWindow.ShowAndRun()
}

func handleFile(_ fyne.Position, uris []fyne.URI) {
	uri := uris[0]
	if len(uris) > 1 {
		dialog.NewInformation("Too many files", "Only reading the first file.", mainWindow).Show()
	}

	file, err := os.Open(uri.Path())
	if err != nil {
		dialog.NewError(err, mainWindow).Show()
		return
	}
	defer file.Close()

	riffChunk, err := riff.DeserializerRiff(file)

	if err != nil {
		dialog.NewError(err, mainWindow).Show()
		return
	}

	wave, ok := riffChunk.(*wave.WaveFile)
	if !ok {
		dialog.NewError(ErrExpectedWave, mainWindow).Show()
	} else {
		handleWave(wave)
	}
}

func handleWave(wave *wave.WaveFile) {
	processingDialog := dialog.NewInformation("Processing", "Let me cook...", mainWindow)
	processingDialog.Show()

	go func() {
		err := wave.MapAllSamples(mapSound)
		if err != nil {
			dialog.NewError(err, mainWindow).Show()
			return
		}

		fyne.Do(func() {
			processingDialog.Dismiss()
			dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.NewError(err, mainWindow).Show()
					return
				}

				err = wave.Serialize(writer)

				if err != nil {
					dialog.NewError(err, mainWindow).Show()
					return
				}
			}, mainWindow).Show()
		})
	}()
}
