package poodle

// #include <linux/fs.h>
import "C"

import (
	"golang.org/x/sys/unix"
)

// RereadPartitionTable reloads the partition table for the block
// device identified by file descriptor fd.
func RereadPartitionTable(fd uintptr) error {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, C.BLKRRPART, 0)
	return err
}
