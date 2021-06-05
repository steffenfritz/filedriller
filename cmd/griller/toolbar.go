package main

import (
	"log"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func createToolbar() *widget.Toolbar {
	toolbar := widget.NewToolbar(widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)

	return toolbar
}