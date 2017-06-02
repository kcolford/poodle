package main

import (
	"flag"
	"log"
	"os"
)

var MainHooks = make([]func() error, 0, 32)

func MainHook(fn func() error) interface{} {
	MainHooks = append(MainHooks, fn)
	return nil
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		// no flags is incorrect flags
		flag.Usage()
		os.Exit(2)
	}
	for _, v := range MainHooks {
		err := v()
		if err != nil {
			log.Fatal(err)
		}
	}
}
