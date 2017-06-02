package main

import (
	"flag"
)

var MainHooks []func()error

func MainHook(fn func() error) interface{} {
	MainHooks = append(MainHooks, fn)
	return nil
}

var help = flag.Bool("help", false, "print this help message")

func main() {
	flag.Parse()
	if flag.NArg() == 0 || *help {
		flag.Usage()
		return
	}
	for _, v := range MainHooks {
		err := v()
		if err != nil {
			panic(err)
		}
	}
}
