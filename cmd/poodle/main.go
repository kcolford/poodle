// Run all the commands together.
package main

import (
	"flag"
	"log"

	"github.com/kcolford/poodle"
)

func run() (err error) {
	kalidevice := flag.String("kali-dev", "", "block device to add encrypted kali linux encrypted partition for")
	randsigid := flag.String("rand-sig", "", "identifier to place the random signature in")
	pkg := flag.String("pkg", "main", "package to place generated code into")
	sshaddr := flag.String("ssh", "", "`host:port` specification for a listening SSH server")
	flag.Parse()
	switch {
	case *kalidevice != "":
		err = poodle.KaliAddEncryptedPartition(*kalidevice)
	case *randsigid != "":
		err = poodle.GenerateRandomSignature(*randsigid, *pkg)
	case *sshaddr != "":
		err = poodle.SshServer(*sshaddr)
	default:
		flag.Usage()
	}
	return
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}
