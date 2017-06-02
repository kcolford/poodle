package main

import (
	"flag"
	"os"
)

var partitionedBlkDev = flag.String("blk-dev", "", "re-read the partition table for `blockdev`")
var _ = MainHook(func() error {
	if *partitionedBlkDev == "" {
		return nil
	}
	f, err := os.Open(*partitionedBlkDev)
	if err != nil {
		return err
	}
	defer f.Close()
	return RereadPartitionTable(f.Fd())
})
