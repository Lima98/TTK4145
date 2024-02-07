package main

import (
	"fmt"
	"net"
	"os/exec"
	"time"
)

const proto, addr = "udp", ":20022"

func main() {
	
	fmt.Println("Starting A-program in 2 seconds")
	time.Sleep(2 * time.Second)
	exec.Command("gnome-terminal", "--", "go", "run", "A.go").Run()
	
	for {
		//send noe greier
		send()
		fmt.Println("Sent \"hello\" on port 20022")
		time.Sleep(1 * time.Second)
	}
}

func send(){
	addr2 := "localhost:20022"
	conn, _ := net.Dial(proto, addr2)                                                                                  
	conn.Write([]byte("hello\n"))                                                                               
	  
}


