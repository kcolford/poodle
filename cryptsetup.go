package poodle

import (
	"os/exec"
)

// CryptsetupFormat formats a disk drive to be a luks encrypted
// volume.
func CryptsetupFormat(disk string) error {
	return exec.Command("cryptsetup", "luksFormat", disk).Run()
}
