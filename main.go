package main

import (
	"fmt"

	"github.com/sakti/filebroker/listener/sshd"
)

func sayHello(obj string) string {
	return fmt.Sprintf("Hello %s", obj)
}

func main() {
	fmt.Println(sayHello("file broker"))
	sshd.Listen(2000)
}
