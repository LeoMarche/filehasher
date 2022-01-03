package main

import (
	"flag"
	"fmt"
	"path"
	"path/filepath"

	"github.com/LeoMarche/filehasher/pkg/copyutils"
)

func main() {
	var src = flag.String("s", "-1", "the source directory")
	var dst = flag.String("d", "-1", "the destination directory")
	var tries = flag.Int("n", 5, "number of tries when copying")

	flag.Parse()

	foldersrc := path.Base(*src)
	folderdst := path.Base(*dst)

	if foldersrc != folderdst {
		*dst = filepath.Join(*dst, foldersrc)
	}

	err := copyutils.CopyTree(*src, *dst, *tries)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Copy terminated with success !")
	}
}
