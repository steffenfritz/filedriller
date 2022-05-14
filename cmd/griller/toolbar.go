package main

import (
	"log"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func createToolbar() *widget.Toolbar {
	toolbar := widget.NewToolbar(widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)

	return toolbar
}
