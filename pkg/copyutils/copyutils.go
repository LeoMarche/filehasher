package copyutils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/LeoMarche/filehasher/pkg/hasherutils"
	"github.com/cespare/xxhash"
)

// copyFile copies a source file to multiple destination files
func copyFile(src io.Reader, dst []*os.File, sizeCopied *int64) error {

	var err error
	var ew error
	var nw int

	// Make copy buffer
	size := 32 * 1024
	if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
		if l.N < 1 {
			size = 1
		} else {
			size = int(l.N)
		}
	}
	buf := make([]byte, size)

	// Execute the copy
	for {

		// Reads the src file
		nr, er := src.Read(buf)
		if nr > 0 {

			// Updates for progressBar
			*sizeCopied += int64(nr)

			// Write to all destination files
			for _, d := range dst {
				nw, ew = d.Write(buf[0:nr])
				if nw < 0 || nr < nw {
					nw = 0
					if ew == nil {
						ew = errors.New("invalid write result")
					}
				}
				if ew != nil {
					break
				}
			}

			// Check errors during write
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return err
}

// safeCopyFile copies a file and verify that the source hash is the same as the destination hash
func safeCopyFile(src string, dst []string, retries int, sizeCopied *int64) error {

	// Variables
	var validCopy bool = false
	var err, err2 error
	var sourceFileStat fs.FileInfo
	var hashsrc, hashdst []byte

	// Number of retries
	t := 0

	// While the copy isn't conform and we didn't exceed the number of retries
	for !validCopy && t < retries {

		// Checks that the source file is correct
		sourceFileStat, err = os.Stat(src)
		if err != nil {
			return err
		}
		if !sourceFileStat.Mode().IsRegular() {
			return fmt.Errorf("%s is not a regular file", src)
		}

		// Opens the source file
		source, err := os.Open(src)
		if err != nil {
			return err
		}
		defer source.Close()

		var destination []*os.File

		for _, f := range dst {

			// Opens the destination file
			d, err := os.Create(f)
			if err != nil {
				return err
			}
			destination = append(destination, d)
			defer d.Close()

		}

		// Executes the copy
		err = copyFile(source, destination, sizeCopied)
		if err != nil {
			return err
		}

		// Hashes the source and the destination
		h := xxhash.New()
		h2 := xxhash.New()
		hashsrc, err = hasherutils.HashFile(src, h)

		validCopy = true
		for _, f := range dst {
			hashdst, err2 = hasherutils.HashFile(f, h2)

			// Checks that the files are the same
			if err == nil && err2 == nil {
				if bytes.Equal(hashsrc, hashdst) {
					validCopy = true
				}
			}
		}

		t++
	}

	// If we exceeded the maximum number of retries
	if t >= retries {
		return fmt.Errorf("number of retries exceeded on file %s, please verify manually", src)
	}

	return nil
}

// SafeCopyTree copies a whole directory and its content
// It verifies that all files copied have the same
// source and destination hashes
func SafeCopyTree(src string, dst []string, retries int, sizeCopied *int64) error {

	// Variables
	var newsrcpath string
	var newdstpath []string
	var dr bool
	var err error
	var fileInfo []fs.FileInfo

	// Verify that src is a directory
	dr, err = hasherutils.IsDirectory(src)
	if !dr {
		return fmt.Errorf("%s is not a directory", src)
	}
	if err != nil {
		return err
	}

	// Gets infos on the directory and verify
	fileInfo, err = ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range dst {
		// Creates the destination directory
		err = os.MkdirAll(f, fs.ModePerm)
		if err != nil {
			return err
		}
	}

	// Recursively executes copy on the files and subdirectories
	for _, file := range fileInfo {
		newdstpath = newdstpath[:0]
		newsrcpath = filepath.Join(src, file.Name())
		for _, f := range dst {
			newdstpath = append(newdstpath, filepath.Join(f, file.Name()))
		}
		dr, err = hasherutils.IsDirectory(newsrcpath)
		if err != nil {
			return err
		}
		if dr {
			err = SafeCopyTree(newsrcpath, newdstpath, retries, sizeCopied)
		} else {
			err = safeCopyFile(newsrcpath, newdstpath, retries, sizeCopied)
		}
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
