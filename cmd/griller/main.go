package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strings"

	fdr "github.com/dla-marbach/filedriller"
)

// we make these two globally available for usage inside imported customized loggers
var logfilevalue widget.Entry
var elogfilevalue widget.Entry

// Version holds the version of filedriller
var Version string
// Build holds the sha1 fingerprint of the build
var Build string
// SigFile holds the download date of the signature file
var SigFile string

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
	if len(outputfilevalue.Text) == 0 {
		outputfilevalue.Text = "info.csv"
	}

	// Log output
	logfilevalue, logfilefield := genericInput("Log file")
	if len(logfilevalue.Text) == 0 {
		logfilevalue.Text = "logs.txt"
	}

	// Error log output
	elogfilevalue, elogfilefield := genericInput("Error log file")
	if len(elogfilevalue.Text) == 0 {
		elogfilevalue.Text = "errorlogs.txt"
	}

	// Fixity 
	fixitywidget := widget.NewSelect([]string{"md5", "sha1", "sha256", "sha512", "blake2b-512"}, func(value string) {
		conf.HashAlg = value
	})
	fixitywidget.PlaceHolder = "Fixity"

	mainInputContainer := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		rootdirfield,
		outputfilefield,
		logfilefield,
		elogfilefield)

	okbutton := widget.NewButton("Start", func() {
		fdr.CreateLogger(logfilevalue.Text)
		fdr.CreateErrorLogger(elogfilevalue.Text)
		if !strings.HasSuffix(rootdirvalue.Text, "/") {
			rootdirvalue.Text = rootdirvalue.Text + "/"
		}
		filelist := fdr.CreateFileList(rootdirvalue.Text)
		progressbar.Max = float64(len(filelist))
		conf.OFile = outputfilevalue.Text
		resultList := fdr.IdentifyFilesGUI(filelist, false, conf, &progressbar.Max)
		fdr.WriteCSV(&outputfilevalue.Text, &conf.HashAlg, resultList)
		fdr.WriteLogfile(Version, Build, SigFile,conf.HashAlg,false, conf.Entro, filelist, resultList)
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
