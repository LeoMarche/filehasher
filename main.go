package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/LeoMarche/filehasher/pkg/copyutils"
	"github.com/gotk3/gotk3/gtk"
)

var WINDOW_HEIGHT int = 600
var WINDOW_WIDTH int = 400

func choose(l *gtk.Entry) {
	dlg, _ := gtk.FileChooserDialogNewWith2Buttons(
		"Choose source folder", nil, gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
		"Choose Folder", gtk.RESPONSE_OK, "Cancel", gtk.RESPONSE_CANCEL,
	)
	dlg.SetDefaultResponse(gtk.RESPONSE_OK)
	response := dlg.Run()
	if response == gtk.RESPONSE_OK {
		filename := dlg.GetFilename()
		l.SetText(filename)
	}
	dlg.Destroy()
}

func verifyFolder(srcPath, dstPath string) error {
	if srcPath == "" || dstPath == "" {
		return fmt.Errorf("at least one directory is empty string")
	}
	s1, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	if !s1.IsDir() {
		return fmt.Errorf("specified string is not a directory")
	}
	s2, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	if !s2.IsDir() {
		return fmt.Errorf("specified string is not a directory")
	}

	return nil
}

func setProgressBar(pB *gtk.ProgressBar, startSync *gtk.Button, percent float64) {
	pB.SetFraction(percent)
	if percent != 1.0 {
		startSync.SetSensitive(false)
	}
}

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

func startCopy(src, dst string, retries int, pB *gtk.ProgressBar, startSync *gtk.Button) {
	ds, err := dirSize(src)
	if err != nil {
		startSync.SetLabel("Failed !")
	}
	var pe int64 = 0
	go func() {
		err := copyutils.CopyTree(src, dst, retries, &pe)
		if err != nil {
			startSync.SetLabel("Failed !")
		}
	}()
	for pe != ds {
		setProgressBar(pB, startSync, float64(pe)/float64(ds))
		time.Sleep(1 * time.Second)
	}
	setProgressBar(pB, startSync, float64(pe)/float64(ds))
	startSync.SetLabel("Finished !")
}

func main() {

	// New GTK window
	gtk.Init(nil)
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("filehasher")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// New Grid
	grid, _ := gtk.GridNew()
	win.Add(grid)

	// Add button for Source Folder
	sourceFolder, _ := gtk.EntryNew()
	sourceFolder.SetMaxWidthChars(100)
	grid.Add(sourceFolder)
	chooseSourceFolder, _ := gtk.ButtonNewWithLabel("Source Folder")
	grid.Add(chooseSourceFolder)

	// Add button for destination folder
	destinationFolder, _ := gtk.EntryNew()
	grid.Attach(destinationFolder, 0, 1, 1, 1)
	chooseDestinationFolder, _ := gtk.ButtonNewWithLabel("Destination Folder")
	grid.Attach(chooseDestinationFolder, 1, 1, 1, 1)

	// Add button for starting sync
	startSync, _ := gtk.ButtonNewWithLabel("Start Syncing")
	grid.Attach(startSync, 0, 2, 2, 1)

	// Add progress bar
	pBar, _ := gtk.ProgressBarNew()
	grid.Attach(pBar, 0, 3, 2, 1)

	pBar.SetFraction(0.0)

	chooseSourceFolder.Connect("clicked", func() {
		choose(sourceFolder)
		t1, err1 := sourceFolder.GetText()
		t2, err2 := destinationFolder.GetText()
		if err1 == nil && err2 == nil && verifyFolder(t1, t2) == nil {
			startSync.Show()
		}
	})

	chooseDestinationFolder.Connect("clicked", func() {
		choose(destinationFolder)
		t1, err1 := sourceFolder.GetText()
		t2, err2 := destinationFolder.GetText()
		if err1 == nil && err2 == nil && verifyFolder(t1, t2) == nil {
			startSync.Show()
		} else {
			startSync.Hide()
		}
	})

	startSync.Connect("clicked", func() {
		t1, err1 := sourceFolder.GetText()
		t2, err2 := destinationFolder.GetText()
		if err1 != nil || err2 != nil {
			fmt.Println("Error with source and destination folders")
		}
		go startCopy(t1, t2, 5, pBar, startSync)
		startSync.SetSensitive(false)
	})

	// Show window
	win.SetDefaultSize(grid.GetSizeRequest())
	win.ShowAll()
	startSync.Hide()
	gtk.Main()
}
