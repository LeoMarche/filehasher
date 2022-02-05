package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/LeoMarche/filehasher/pkg/copyutils"
	"github.com/therecipe/qt/widgets"
)

// dirSize returns the cumulated size of
// all files in a directory and its sub directories
func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// isEmpty returns wether or not a directory
// is empty
func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// startCopy starts the checksumed copy using goroutines
// This way, GUI is still ressponsive during copy
func startCopy(src string, dst []*widgets.QLabel, retries int, pB *widgets.QProgressBar, startSync *widgets.QPushButton) {

	var failed bool = false

	// Deactivate copy button
	startSync.SetEnabled(false)
	startSync.SetText("Copy in progress ...")

	// Setup progress bar
	ds, err := dirSize(src)
	pB.SetMaximum(int(ds))
	if err != nil {
		startSync.SetText("Failed !")
	}

	// Run the checksumed copy in a separated goroutine
	var pe int64 = 0

	var dstString []string

	for _, f := range dst {
		dstString = append(dstString, f.Text())
	}

	go func() {
		err := copyutils.SafeCopyTree(src, dstString, retries, &pe)
		if err != nil {
			startSync.SetText("Failed !")
			failed = true
		}
	}()

	// Actuate the progress bar once every 500ms
	for pe != ds {
		pB.SetValue(int(pe))
		time.Sleep(500 * time.Millisecond)
		if failed {
			return
		}
	}

	// Actuate the progress bar once at 100%
	pB.SetValue(int(pe))

	// Set button text
	startSync.SetText("Finished !")
}

// selectDirectory brings up a Select Directory dialog
// and update GUI accordingly with the selected directory
func selectDirectory(label *widgets.QLabel) {
	path := widgets.QFileDialog_GetExistingDirectory(nil, "Choose directory", "", widgets.QFileDialog__ShowDirsOnly)
	if path != "" {
		label.SetText(path)
	} else {
		fmt.Println("not exists")
	}
}

// checkDirs checks if the two directories selected
// are eligible for copy
func checkDirs(src string, dst []*widgets.QLabel, but *widgets.QPushButton) {

	eD := true
	eL := false

	for _, d := range dst {
		// Checks that destination folders are empty
		emptyDestination, _ := isEmpty(d.Text())
		if !emptyDestination && d.Text() != "" {
			eD = false
			but.SetText("Destination isn't empty !")
			but.SetStyleSheet("QPushButton {color: red;}")
		}
		if d.Text() == "" {
			eL = true
		}
	}

	// Check that src and dst are correct
	if src != "" && !eL && eD {
		but.SetStyleSheet("")
		but.SetText("Start copy !")
		but.SetEnabled(true)
	}
}

// Main loop
func main() {

	var DestLabels []*widgets.QLabel

	// Setup GUI
	widgets.NewQApplication(len(os.Args), os.Args)

	// Setup source box
	var (
		SourceFileChooser = widgets.NewQPushButton2("Choose source folder !", nil)
		SourceGroup       = widgets.NewQGroupBox2("Source", nil)
		SourceLabel       = widgets.NewQLabel2("", nil, 0)
	)

	// Setup destination box
	var (
		DestFileChooser    = widgets.NewQPushButton2("Choose destination folder !", nil)
		DestGroup          = widgets.NewQGroupBox2("Destination", nil)
		DestLabel          = widgets.NewQLabel2("", nil, 0)
		DestAddFileChooser = widgets.NewQPushButton2("Add a destination folder !", nil)
	)

	DestLabels = append(DestLabels, DestLabel)

	// Setup start syncing box
	var (
		ProgressBar = widgets.NewQProgressBar(nil)
		StartGroup  = widgets.NewQGroupBox2("Start", nil)
		StartButton = widgets.NewQPushButton2("Start copy !", nil)
	)

	// Set initial state of "start syncing box"
	ProgressBar.SetMinimum(0)
	ProgressBar.SetMaximum(100)
	ProgressBar.SetValue(0)
	StartButton.SetEnabled(false)

	// Connect buttons to their functions
	SourceFileChooser.ConnectClicked(func(checked bool) {
		selectDirectory(SourceLabel)
		checkDirs(SourceLabel.Text(), DestLabels, StartButton)

	})
	DestFileChooser.ConnectClicked(func(checked bool) {
		selectDirectory(DestLabel)
		checkDirs(SourceLabel.Text(), DestLabels, StartButton)
	})
	StartButton.ConnectClicked(func(checked bool) {
		StartButton.SetEnabled(false)
		go startCopy(SourceLabel.Text(), DestLabels, 5, ProgressBar, StartButton)
	})

	// Setup window Layout
	var SourceLayout = widgets.NewQGridLayout2()
	SourceLayout.AddWidget(SourceFileChooser)
	SourceLayout.AddWidget(SourceLabel)
	SourceGroup.SetLayout(SourceLayout)

	var DestLayout = widgets.NewQGridLayout2()
	DestLayout.AddWidget(DestFileChooser)
	DestLayout.AddWidget(DestLabel)
	DestGroup.SetLayout(DestLayout)

	var StartLayout = widgets.NewQGridLayout2()
	StartLayout.AddWidget(ProgressBar)
	StartLayout.AddWidget(StartButton)
	StartGroup.SetLayout(StartLayout)

	var layout = widgets.NewQGridLayout2()
	layout.AddWidget(SourceGroup)
	layout.AddWidget(DestGroup)
	layout.AddWidget(DestAddFileChooser)
	layout.AddWidget(StartGroup)

	DestAddFileChooser.ConnectClicked(func(checked bool) {
		newDestFileChooser := widgets.NewQPushButton2("Choose destination folder !", nil)
		newDestLabel := widgets.NewQLabel2("", nil, 0)
		DestLabels = append(DestLabels, newDestLabel)
		newDestFileChooser.ConnectClicked(func(checked bool) {
			selectDirectory(newDestLabel)
			checkDirs(SourceLabel.Text(), DestLabels, StartButton)
		})
		DestLayout.AddWidget(newDestFileChooser)
		DestLayout.AddWidget(newDestLabel)
	})

	// Setup window
	var window = widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Filehasher")

	var centralWidget = widgets.NewQWidget(window, 0)
	centralWidget.SetLayout(layout)
	window.SetCentralWidget(centralWidget)

	window.Show()

	// start GUI
	widgets.QApplication_Exec()
}
