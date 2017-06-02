package poodle

import (
	"errors"
	"fmt"
	"os"

	"github.com/rekby/mbr"
)

const kaliEnctyptedPartno = 3
const kaliEncryptedLocation = 0x100000000 // 4G into the device
const sectorSize = 512

func kaliAddPartition(disk string) error {
	// open the disk and check the partition table
	f, err := os.Open(disk)
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
	info, err := os.Stat(disk)
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
	f, err = os.OpenFile(disk, os.O_WRONLY, 0644)
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

// GetPartitionName returns the name of the block device that
// corresponds to the partition numbered partno on the enclosing block
// device blkdev.
func GetPartitionName(blkdev string, partno int) string {
	lastchar := blkdev[len(blkdev)-1]
	if '0' <= lastchar || lastchar <= '9' {
		blkdev += "p"
	}
	return blkdev + fmt.Sprint(partno)
}

// KaliAddEncryptedPartition adds an encrypted partition to the block
// device disk that is supposed to have the Kali Linux ISO file imaged
// to it.
func KaliAddEncryptedPartition(disk string) error {
	// add the partition
	err := kaliAddPartition(disk)
	if err != nil {
		return err
	}

	// run cryptsetup to format the partition
	partdev := GetPartitionName(disk, kaliEnctyptedPartno)
	err = CryptsetupFormat(partdev)
	return err
}
