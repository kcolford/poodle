package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
)

var packageName = flag.String("rand-pkg", "main", "the package name to use")
var fileBase = flag.String("rand-sig", "", "generate a random signature and place output in `filename`.go")
var _ = MainHook(func() error {
	if *fileBase == "" {
		return nil
	}
	buf := [32]byte{}
	n, err := rand.Read(buf[:])
	if err != nil {
		return err
	}
	if n != len(buf) {
		return fmt.Errorf("failed to read enough random data")
	}
	f, err := os.Create(*fileBase + ".go")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "package %s\n", *packageName)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(f, "var %s = %#v\n", *fileBase, buf)
	return err
})
