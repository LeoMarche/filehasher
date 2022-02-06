package hasherutils

import (
	"hash"
	"path/filepath"
	"testing"

	"github.com/cespare/xxhash"
	"github.com/stretchr/testify/assert"
)

type testFileHash struct {
	path          string
	expectedHash  []byte
	expectedError bool
	hasher        hash.Hash
}

type testIsDirectory struct {
	path          string
	result        bool
	expectedError bool
}

func TestHashFile(t *testing.T) {
	basepath := filepath.Join("..", "..", "tests", "testHashFile")
	t1 := testFileHash{
		path:          "text",
		expectedHash:  []byte{0xb3, 0xe7, 0xce, 0x50, 0x6e, 0xd6, 0xfd, 0xfd},
		expectedError: false,
		hasher:        xxhash.New(),
	}
	ts := []testFileHash{t1}

	for _, test := range ts {
		hash, err := HashFile(filepath.Join(basepath, test.path), test.hasher)
		if test.expectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedHash, hash)
		}
	}

}

func TestIsDirectory(t *testing.T) {
	basepath := filepath.Join("..", "..", "tests", "testIsDirectory")

	t1 := testIsDirectory{
		path:          "dir",
		result:        true,
		expectedError: false,
	}

	t2 := testIsDirectory{
		path:          "file",
		result:        false,
		expectedError: false,
	}

	t3 := testIsDirectory{
		path:          "notexist",
		result:        false,
		expectedError: true,
	}

	ts := []testIsDirectory{t1, t2, t3}

	for _, test := range ts {
		re, err := IsDirectory(filepath.Join(basepath, test.path))
		if test.expectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, test.result, re)
	}
}
