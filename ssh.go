package poodle

import (
	"fmt"
	"net"

	"crypto/rand"
	"crypto/rsa"
	"errors"
	"io"

	"golang.org/x/crypto/ssh"
	"os"
)

func SshServer(laddr string) error {
	// setup ssh server config
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if string(pass) == "pass" {
				return nil, nil
			}
			return nil, errors.New("password rejected")
		},
	}
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		return err
	}
	config.AddHostKey(signer)

	// setup socket listener
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		con, newchannels, reqs, err := ssh.NewServerConn(conn, config)
		if err != nil {
			return err
			//continue
		}

		// discard out of band requests
		go func() {
			for req := range reqs {
				req.Reply(false, nil)
			}
		}()

		// respond to the incoming multiplexed channels
		go func() {
			defer con.Close()
			for newchannel := range newchannels {
				fmt.Printf("%+v\n", newchannel) // output for debug
				channel, reqs, err := newchannel.Accept()
				if err != nil {
					continue
				}

				go func() {
					for req := range reqs {
						fmt.Printf("%s %s\n", req.Type, req.Payload)
						req.Reply(true, nil)
					}
				}()

				go func() {
					defer channel.Close()
					channel.CloseWrite()
					io.Copy(os.Stdout, channel)
				}()
			}
		}()
	}

	return nil
}
