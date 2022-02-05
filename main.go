package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/LeoMarche/filehasher/pkg/copyutils"
	"github.com/therecipe/qt/widgets"
)

var WINDOW_HEIGHT int = 600
var WINDOW_WIDTH int = 400

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

func startCopy(src, dst string, retries int, pB *widgets.QProgressBar, startSync *widgets.QPushButton) {
	ds, err := dirSize(src)
	pB.SetMaximum(int(ds))
	if err != nil {
		startSync.SetText("Failed !")
	}
	var pe int64 = 0
	go func() {
		err := copyutils.CopyTree(src, dst, retries, &pe)
		if err != nil {
			startSync.SetText("Failed !")
		}
	}()
	for pe != ds {
		pB.SetValue(int(pe))
		time.Sleep(500 * time.Millisecond)
	}
	pB.SetValue(int(pe))
	startSync.SetText("Finished !")
}

func selectDirectory(label *widgets.QLabel) {
	path := widgets.QFileDialog_GetExistingDirectory(nil, "Choose directory", "", widgets.QFileDialog__ShowDirsOnly)
	if path != "" {
		label.SetText(path)
	} else {
		fmt.Println("not exists")
	}
}

func checkDirs(d1, d2 string, but *widgets.QPushButton) {
	if d1 != "" && d2 != "" {
		but.SetEnabled(true)
	}
}

func main() {
	widgets.NewQApplication(len(os.Args), os.Args)

	var (
		echoSourceFileChooser = widgets.NewQPushButton2("Choose source folder !", nil)
		echoSourceGroup       = widgets.NewQGroupBox2("Source", nil)
		echoSourceLabel       = widgets.NewQLabel2("", nil, 0)
	)

	var (
		echoDestFileChooser = widgets.NewQPushButton2("Choose destination folder !", nil)
		echoDestGroup       = widgets.NewQGroupBox2("Destination", nil)
		echoDestLabel       = widgets.NewQLabel2("", nil, 0)
	)

	var (
		echoProgressBar = widgets.NewQProgressBar(nil)
		echoStartGroup  = widgets.NewQGroupBox2("Start", nil)
		echoStartButton = widgets.NewQPushButton2("Start copy !", nil)
	)

	echoProgressBar.SetMinimum(0)
	echoProgressBar.SetMaximum(100)
	echoProgressBar.SetValue(0)
	echoStartButton.SetEnabled(false)

	echoSourceFileChooser.ConnectClicked(func(checked bool) {
		selectDirectory(echoSourceLabel)
		checkDirs(echoDestLabel.Text(), echoSourceLabel.Text(), echoStartButton)

	})
	echoDestFileChooser.ConnectClicked(func(checked bool) {
		selectDirectory(echoDestLabel)
		checkDirs(echoDestLabel.Text(), echoSourceLabel.Text(), echoStartButton)
	})
	echoStartButton.ConnectClicked(func(checked bool) {
		echoStartButton.SetEnabled(false)
		go startCopy(echoSourceLabel.Text(), echoDestLabel.Text(), 5, echoProgressBar, echoStartButton)
	})

	var echoSourceLayout = widgets.NewQGridLayout2()
	echoSourceLayout.AddWidget2(echoSourceFileChooser, 0, 0, 0)
	echoSourceLayout.AddWidget2(echoSourceLabel, 1, 0, 0)
	echoSourceGroup.SetLayout(echoSourceLayout)

	var echoDestLayout = widgets.NewQGridLayout2()
	echoDestLayout.AddWidget2(echoDestFileChooser, 0, 0, 0)
	echoDestLayout.AddWidget2(echoDestLabel, 1, 0, 0)
	echoDestGroup.SetLayout(echoDestLayout)

	var echoStartLayout = widgets.NewQGridLayout2()
	echoStartLayout.AddWidget2(echoProgressBar, 0, 0, 0)
	echoStartLayout.AddWidget2(echoStartButton, 1, 0, 0)
	echoStartGroup.SetLayout(echoStartLayout)

	var layout = widgets.NewQGridLayout2()
	layout.AddWidget2(echoSourceGroup, 0, 0, 0)
	layout.AddWidget2(echoDestGroup, 1, 0, 0)
	layout.AddWidget2(echoStartGroup, 2, 0, 0)

	var window = widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Filehasher")

	var centralWidget = widgets.NewQWidget(window, 0)
	centralWidget.SetLayout(layout)
	window.SetCentralWidget(centralWidget)

	window.Show()

	widgets.QApplication_Exec()
}
