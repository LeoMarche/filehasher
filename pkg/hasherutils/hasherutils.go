package hasherutils

import (
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func HashFile(path string, h hash.Hash) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h.Sum((nil)), nil
}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func HashFolder(path string, h hash.Hash) ([]byte, error) {
	var newpath string
	var dr bool
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range fileInfo {
		newpath = filepath.Join(path, file.Name())
		dr, err = IsDirectory(newpath)
		if err != nil {
			return nil, err
		}
		if dr {
			_, err = HashFolder(newpath, h)
		} else {
			_, err = HashFile(newpath, h)
		}
		if err != nil {
			return nil, err
		}
	}

	return h.Sum((nil)), nil
}
