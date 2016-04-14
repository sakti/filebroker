package main

import (
	"fmt"

	"github.com/sakti/filebroker/listener/sshd"
)

func main() {
	fmt.Println("file broker")
	sshd.Listen(2000)
}
