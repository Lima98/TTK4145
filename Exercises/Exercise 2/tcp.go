package main

import (
	"fmt"
	"net"
	"time"
)

func main(){
	go tcpCommunicate()
	go acceptCommunicate()

	time.Sleep(5 * time.Second)
}

func tcpCommunicate(){
	const proto, addr = "tcp", "10.100.23.129:33546"
	
	sock, _ := net.Dial(proto, addr)
	sock.Write([]byte("Connect to: 10.100.23.32:33546\x00"))

	for {
		sock.Write([]byte("Kobling: Oss -> SERVER \x00"))
		buf := make([]byte, 1024)
		n, _ := sock.Read(buf)

		fmt.Printf("Received %d bytes: %s\n", n, string(buf[:n]))
		time.Sleep(1 * time.Second)
	}
}

func acceptCommunicate(){
	const proto, addr = "tcp", "10.100.23.129:33546"
	
	sock, _ := net.Listen(proto, ":33546")
	conn, _ := sock.Accept()
	for{
		conn.Write([]byte("Kobling: SERVER -> OSS \x00"))
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)

		fmt.Printf("Received %d bytes: %s\n", n, string(buf[:n]))
		time.Sleep(1 * time.Second)
	}

}

