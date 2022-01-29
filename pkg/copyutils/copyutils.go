package copyutils

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/LeoMarche/filehasher/pkg/hasherutils"
	"github.com/cespare/xxhash"
)

func copyFile(src, dst string, retries int) (int64, error) {
	validCopy := false
	var nBytes int64
	var err, err2 error
	var sourceFileStat fs.FileInfo
	var hashsrc, hashdst []byte
	t := 0
	for !validCopy && t < retries {
		sourceFileStat, err = os.Stat(src)
		if err != nil {
			return 0, err
		}

		if !sourceFileStat.Mode().IsRegular() {
			return 0, fmt.Errorf("%s is not a regular file", src)
		}

		source, err := os.Open(src)
		if err != nil {
			return 0, err
		}
		defer source.Close()

		destination, err := os.Create(dst)
		if err != nil {
			return 0, err
		}
		defer destination.Close()
		nBytes, err = io.Copy(destination, source)
		if err != nil {
			return 0, err
		}
		h := xxhash.New()
		h2 := xxhash.New()
		hashsrc, err = hasherutils.HashFile(src, h)
		hashdst, err2 = hasherutils.HashFile(dst, h2)

		if err == nil && err2 == nil {
			if bytes.Equal(hashsrc, hashdst) {
				validCopy = true
			}
		}
		t++
	}

	if t >= retries {
		return 0, fmt.Errorf("number of retries exceeded on file %s, please verify manually", src)
	}

	return nBytes, nil
}

func CopyTree(src, dst string, retries int) error {
	var newsrcpath, newdstpath string
	var dr bool
	t := 0
	var err error
	var fileInfo []fs.FileInfo

	fileInfo, err = ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	dr, err = hasherutils.IsDirectory(src)
	if !dr {
		return fmt.Errorf("%s is not a directory", src)
	}
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, fs.ModePerm)
	if err != nil {
		return err
	}

	for _, file := range fileInfo {
		newsrcpath = filepath.Join(src, file.Name())
		newdstpath = filepath.Join(dst, file.Name())
		dr, err = hasherutils.IsDirectory(newsrcpath)
		if err != nil {
			return err
		}
		if dr {
			err = CopyTree(newsrcpath, newdstpath, retries)
		} else {
			_, err = copyFile(newsrcpath, newdstpath, retries)
		}
		if err != nil {
			fmt.Println(err)
		}
	}

	if t >= retries {
		return fmt.Errorf("number of retries exceeded")
	}

	return nil
}
