package main

import (
	"fmt"
	"net"
)

const proto, addr = "udp", "localhost:20022"

func main() {
	
	fmt.Println("I am alive")

	for{
		receive()
	}

}

func receive() {
	conn, _ := net.ListenPacket(proto, addr)
	for {	
		buf := make([]byte, 1024)
		num_of_bytes, source, _ := conn.ReadFrom(buf)
		fmt.Println(string(buf[:num_of_bytes]))
		conn.WriteTo(buf,source)
	}
}
