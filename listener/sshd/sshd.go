package sshd

// see: https://blog.gopheracademy.com/advent-2015/ssh-server-in-go/

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

// Listen start SSH server
func Listen(port int) {
	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			fmt.Println(strings.TrimSpace(string(ssh.MarshalAuthorizedKey(key))))
			return &ssh.Permissions{Extensions: map[string]string{"key-id": "test"}}, nil
		},
	}

	keyPath := filepath.Join("config", "ssh/key.rsa")

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(keyPath), os.ModePerm)
		out, err := exec.Command("ssh-keygen", "-f", keyPath, "-t", "rsa", "-N", "").Output()
		if err != nil {
			panic(fmt.Sprintf("error: %s", err))
		}
		fmt.Println(out)
		fmt.Println("ssh server private key created")
	}

	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic("failed to load ssh server private key")
	}

	key, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		panic("failed to parse ssh server private key")
	}

	config.AddHostKey(key)

	listen(config, port)

}

func listen(config *ssh.ServerConfig, port int) {
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	fmt.Println("Listening on 0.0.0.0:" + strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// handshake
		sConn, chans, reqs, err := ssh.NewServerConn(conn, config)
		if err != nil {
			continue
		}

		go ssh.DiscardRequests(reqs)
		go handleServerConn(sConn.Permissions.Extensions["key-id"], chans)

	}
}

func handleServerConn(keyID string, chans <-chan ssh.NewChannel) {
	for newChan := range chans {
		if newChan.ChannelType() != "session" {
			newChan.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		ch, reqs, err := newChan.Accept()

		if err != nil {
			continue
		}

		go func(in <-chan *ssh.Request) {
			defer ch.Close()
			for req := range in {
				fmt.Println(req.Type)
				fmt.Println(req.Payload)
			}
		}(reqs)
	}

}
