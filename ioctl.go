package main

// #include <linux/fs.h>
// #include <linux/blkpg.h>
import "C"

import "golang.org/x/sys/unix"

func Ioctl(fd, cmd, arg uintptr) error {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, cmd, arg)
	return err
}

func RereadPartitionTable(fd uintptr) error {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, C.BLKRRPART, 0)
	return err
}
