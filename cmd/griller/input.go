package main

import "fyne.io/fyne/v2/widget"

// genericInput creates a generic form widget
func genericInput(formkeytext string) (*widget.Entry, *widget.Form) {
	entry := widget.NewEntry()
	entry.PlaceHolder = "Please enter something"
	FormItem := &widget.FormItem{
		Text:   formkeytext,
		Widget: entry,
	}
	entryfield := widget.NewForm(FormItem)

	return entry, entryfield
}
