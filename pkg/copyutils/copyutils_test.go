package copyutils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cespare/xxhash"
	"github.com/stretchr/testify/assert"
	"github.com/udhos/equalfile"
)

type copyFileTest struct {
	dests         []string
	src           string
	expectedError bool
	success       bool
	expectedSize  int64
}

type copyTreeTest struct {
	dests         []string
	src           string
	expectedError bool
	success       bool
	expectedSize  int64
	files         []string
}

func TestCopyFile(t *testing.T) {

	ct1 := copyFileTest{
		dests:         []string{filepath.Join("dst1", "text"), filepath.Join("dst2", "text")},
		src:           filepath.Join("src", "text"),
		expectedError: false,
		success:       true,
		expectedSize:  24,
	}

	ct2 := copyFileTest{
		dests:         []string{filepath.Join("dst1", "png.png"), filepath.Join("dst2", "png.png")},
		src:           filepath.Join("src", "png.png"),
		expectedError: false,
		success:       true,
		expectedSize:  3844738,
	}

	ct := []copyFileTest{ct1, ct2}

	basepath := filepath.Join("..", "..", "tests", "testCopyFile")

	for _, te := range ct {
		src, _ := os.Open(filepath.Join(basepath, te.src))
		var dests []*os.File
		var sizeCopied int64 = 0

		for _, d := range te.dests {
			os.Remove(filepath.Join(basepath, d))
			dFile, _ := os.Create(filepath.Join(basepath, d))
			dests = append(dests, dFile)
		}

		err := copyFile(src, dests, &sizeCopied)

		if te.expectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		for _, d := range te.dests {
			cmp := equalfile.NewMultiple(nil, equalfile.Options{}, xxhash.New(), true)
			equal, err := cmp.CompareFile(filepath.Join(basepath, te.src), filepath.Join(basepath, d))
			if te.success {
				assert.NoError(t, err)
				assert.True(t, equal)
			} else {
				assert.Error(t, err)
				assert.False(t, equal)
			}
		}

		if te.success {
			assert.Equal(t, te.expectedSize, sizeCopied)
		} else {
			assert.Equal(t, 0, sizeCopied)
		}

		for _, d := range te.dests {
			os.Remove(filepath.Join(basepath, d))
		}
	}
}

func TestSafeCopyFile(t *testing.T) {

	ct1 := copyFileTest{
		dests:         []string{filepath.Join("dst1", "text"), filepath.Join("dst2", "text")},
		src:           filepath.Join("src", "text"),
		expectedError: false,
		success:       true,
		expectedSize:  24,
	}

	ct2 := copyFileTest{
		dests:         []string{filepath.Join("dst1", "png.png"), filepath.Join("dst2", "png.png")},
		src:           filepath.Join("src", "png.png"),
		expectedError: false,
		success:       true,
		expectedSize:  3844738,
	}

	ct3 := copyFileTest{
		dests:         []string{filepath.Join("dst1", "notexist"), filepath.Join("dst2", "notexist")},
		src:           filepath.Join("src", "notexist"),
		expectedError: true,
		success:       false,
		expectedSize:  0,
	}

	ct := []copyFileTest{ct1, ct2, ct3}

	basepath := filepath.Join("..", "..", "tests", "testCopyFile")

	for _, te := range ct {
		var dests []string
		var sizeCopied int64 = 0

		for _, d := range te.dests {
			os.Remove(filepath.Join(basepath, d))
			dests = append(dests, filepath.Join(basepath, d))
		}

		err := safeCopyFile(filepath.Join(basepath, te.src), dests, 5, &sizeCopied)

		if te.expectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		for _, d := range te.dests {
			cmp := equalfile.NewMultiple(nil, equalfile.Options{}, xxhash.New(), true)
			equal, err := cmp.CompareFile(filepath.Join(basepath, te.src), filepath.Join(basepath, d))
			if te.success {
				assert.NoError(t, err)
				assert.True(t, equal)
			} else {
				assert.Error(t, err)
				assert.False(t, equal)
			}
		}

		assert.Equal(t, te.expectedSize, sizeCopied)

		for _, d := range te.dests {
			os.Remove(filepath.Join(basepath, d))
		}
	}
}

func TestSafeCopyTree(t *testing.T) {

	ct1 := copyTreeTest{
		dests:         []string{"dst1", "dst2"},
		src:           "src",
		expectedError: false,
		success:       true,
		expectedSize:  24,
		files:         []string{"othertext", filepath.Join("testdir", "text")},
	}

	ct := []copyTreeTest{ct1}

	basepath := filepath.Join("..", "..", "tests", "testCopyTree")

	for _, te := range ct {
		var dests []string
		var sizeCopied int64 = 0

		for _, d := range te.dests {
			os.Remove(filepath.Join(basepath, d))
			dests = append(dests, filepath.Join(basepath, d))
		}

		err := SafeCopyTree(filepath.Join(basepath, te.src), dests, 5, &sizeCopied, &[]error{})

		if te.expectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		for _, d := range te.dests {
			for _, f := range te.files {
				cmp := equalfile.NewMultiple(nil, equalfile.Options{}, xxhash.New(), true)
				equal, err := cmp.CompareFile(filepath.Join(basepath, te.src, f), filepath.Join(basepath, d, f))
				if te.success {
					assert.NoError(t, err)
					assert.True(t, equal)
				} else {
					assert.Error(t, err)
					assert.False(t, equal)
				}
			}
		}

		assert.Equal(t, te.expectedSize, sizeCopied)

		for _, d := range te.dests {
			dir, _ := ioutil.ReadDir(filepath.Join(basepath, d))
			for _, di := range dir {
				os.RemoveAll(filepath.Join(basepath, d, di.Name()))
			}
		}
	}
}
