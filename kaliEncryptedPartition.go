package main

import (
	"fmt"
	"os"

	"errors"
	"flag"

	"github.com/rekby/mbr"
)

const kaliEnctyptedPartno = 3
const kaliEncryptedLocation = 0x100000000 // 4G into the device
const sectorSize = 512

func kaliAddPartition(kaliDisk *string) error {
	// open the disk and check the partition table
	f, err := os.Open(*kaliDisk)
	if err != nil {
		return err
	}
	defer f.Close()
	tbl, err := mbr.Read(f)
	if err != nil {
		return err
	}
	parts := tbl.GetAllPartitions()

	// choose the partition to change
	target := parts[kaliEnctyptedPartno-1]

	// change the partition
	if !target.IsEmpty() {
		return errors.New("the third partition is not empty")
	}
	info, err := os.Stat(*kaliDisk)
	if err != nil {
		return err
	}
	size := info.Size()
	if size%sectorSize != 0 {
		return errors.New("the block device size is not a multiple of the sector size")
	}
	target.SetLBAStart(kaliEncryptedLocation / sectorSize)
	target.SetLBALen(uint32(size / sectorSize))

	// write the partition table back to the device
	f, err = os.OpenFile(*kaliDisk, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	err = tbl.Write(f)
	if err != nil {
		return err
	}

	// update linux's partitions
	err = RereadPartitionTable(f.Fd())
	if err != nil {
		return err
	}

	return nil
}

func GetPartitionName(blkdev string, partno int) string {
	lastchar := blkdev[len(blkdev)-1]
	if '0' <= lastchar || lastchar <= '9' {
		blkdev += "p"
	}
	return blkdev + fmt.Sprint(partno)
}

var kaliDisk = flag.String("kali-disk", "", "format `disk` to have an extra encrypted partition")
var _ = MainHook(func() error {
	if *kaliDisk == "" {
		return nil
	}

	// add the partition
	err := kaliAddPartition(kaliDisk)
	if err != nil {
		return err
	}

	// run cryptsetup to format the partition
	partdev := GetPartitionName(*kaliDisk, kaliEnctyptedPartno)
	err = CryptsetupFormat(partdev)
	if err != nil {
		return err
	}

	return nil
})
