package main

import (
	"fmt"

	"activate/screens"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.NewWithID("activate.app")

	w := a.NewWindow("ACTivate")
	w.SetMaster()

	w.Resize(fyne.NewSize(390, 422))
	w.SetFixedSize(true)
	w.SetContent(screens.Make(a, w))
	w.ShowAndRun()
	tidyUp("Window Closed")
}

func tidyUp(msg string) {
	fmt.Printf("ACTivate Exited: %s\n", msg)
}
