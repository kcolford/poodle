package main

import (
	"os/exec"
)

func CryptsetupFormat(disk string) error {
	return exec.Command("cryptsetup", "luksFormat", disk).Run()
}
