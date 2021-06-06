package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strings"

	fdr "github.com/dla-marbach/filedriller"
)

func main(){
	var conf fdr.Config

	a := app.New()
	w := a.NewWindow("filedriller - GUI")

	// Toolbar
	toolbar := createToolbar() // ToDo

	// Progress bar
	progressbar := widget.NewProgressBar()

	// Input directory
	rootdirvalue, rootdirfield := genericInput("Root directory")

	// Save output file
	outputfilevalue, outputfilefield := genericInput("Output file")

	// Log output
	// logfilevalue, logfilefield := genericInput("Log file")

	// Error log output
	// elogfilevalue, elogfilefield := genericInput("Error log file")

	// Fixity 
	fixitywidget := widget.NewSelect([]string{"md5", "sha1", "sha256", "sha512", "blake2b-512"}, func(value string) {
		conf.HashAlg = value
	})
	fixitywidget.PlaceHolder = "Fixity"

	mainInputContainer := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		rootdirfield,
		outputfilefield)

	okbutton := widget.NewButton("Start", func() {
		if !strings.HasSuffix(rootdirvalue.Text, "/") {
			rootdirvalue.Text = rootdirvalue.Text + "/"
		}
		filelist := fdr.CreateFileList(rootdirvalue.Text)
		progressbar.Max = float64(len(filelist))
		conf.OFile = outputfilevalue.Text
		resultList := fdr.IdentifyFilesGUI(filelist, false, conf, &progressbar.Max)
		fdr.WriteCSV(&outputfilevalue.Text, &conf.HashAlg, resultList)
	})
	quitbutton := widget.NewButton("Quit", func() { a.Quit() })
	buttoncontainer := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), okbutton, quitbutton)

	w.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		toolbar,
		mainInputContainer,
		fixitywidget,
		progressbar,
		buttoncontainer))
	w.ShowAndRun()
}
