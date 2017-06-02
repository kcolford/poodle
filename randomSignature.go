package poodle

import (
	"fmt"
	"math/rand"
	"os"
	"path"
)

// GenerateRandomSignature generates a random byte sequence and places
// it in a Go file for use in code. The identifier should match the
// filename without the .go extension.
func GenerateRandomSignature(filepath, pkgname string) error {
	_, identifier := path.Split(filepath)
	buf := make([]byte, 32)
	n, err := rand.Read(buf)
	if err != nil {
		return err
	}
	if n != len(buf) {
		return fmt.Errorf("failed to read enough random data")
	}
	f, err := os.Create(filepath + ".go")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "package %s\nvar %s = %#v\n", pkgname, identifier, buf)
	return err
}
