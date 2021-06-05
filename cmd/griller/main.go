package main

import (
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	//"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/dla-marbach/filedriller"
)

func main(){
	a := app.New()
	w := a.NewWindow("filedriller - GUI")

	// Toolbar
	toolbar := createToolbar() // ToDo

	// Progress bar
	infinite := widget.NewProgressBarInfinite()
	infinite.Hidden = true
}
