package main

import (
	"crypto/sha256"
	"fmt"

	"github.com/LeoMarche/filehasher/pkg/hasherutils"
)

func main() {
	h := sha256.New()
	fmt.Println(hasherutils.HashFolder("./", h))
}
