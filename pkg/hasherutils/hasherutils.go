package hasherutils

import (
	"hash"
	"io"
	"os"
)

// HashFile computes the checksum of a file using
// the provided hash function
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

// IsDirectory checks wether a path points to a directory or not
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}
