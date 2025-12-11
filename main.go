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
	processingDialog := dialog.NewInformation("Processing", "Working...", mainWindow)
	processingDialog.Show()

	go func() {
		effect(wave, 40*4, 16)
		effect(wave, 48*4, 8)
		effect(wave, 58*4, 16)
		effect(wave, 70*4, 24)
		effect(wave, 88*4, 8)

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

func effect(wave *wave.WaveFile, startBeat, beatLength float64) {
	secondsPerBeat := 60.0 / 125.0 // 125 bpm
	err := applyEffect(wave, startBeat*secondsPerBeat, (startBeat+beatLength)*secondsPerBeat, secondsPerBeat)

	if err != nil {
		fyne.Do(func() {
			dialog.NewError(err, mainWindow).Show()
		})
		return
	}
}
