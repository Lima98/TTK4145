package main

import (
	"fmt"
	"net"
	"time"
)

const proto, addr = "udp", "localhost:20022"

func main() {
	fmt.Println("I am alive")
	go receive()

	for{
		
	}

}
