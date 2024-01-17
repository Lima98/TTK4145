package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

const proto, addr = "udp", ":30000"

func main() {
	
	go receive()

	time.Sleep(1 * time.Second)

	go send()

	time.Sleep(1 * time.Second)
}

func receive() {
	conn, _ := net.ListenPacket(proto, addr)
	buf := make([]byte, 1024)
	num_of_bytes, source, _ := conn.ReadFrom(buf)
	fmt.Println(string(buf[:num_of_bytes]))
	conn.WriteTo(buf,source)
}

func send(){
	addr2 := "#.#.#.255:20022"
	conn, _ := net.Dial(proto, addr2)                                                                                  
	conn.Write([]byte("hello\n"))                                                                               
	buf, _, _ := bufio.NewReader(conn).ReadLine()                                                                     
	fmt.Println("clnt recv", string(buf))  
}