package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const proto, addr = "udp", ":20022"

func main() {
	var counter = 0
	fmt.Println("Starting Backup")
	exec.Command("gnome-terminal", "--", "go", "run", "A.go").Run()
	
	for {
		writeToFile(strconv.FormatInt(int64(counter),10))
		send()
		time.Sleep(1 * time.Second)
		counter++
	}
}

func send(){
	addr2 := "localhost:20022"
	conn, _ := net.Dial(proto,addr2)                                                                                  
	conn.Write([]byte("I'm alive!"))                                                                               
}

func writeToFile(data string){
	os.WriteFile("./data.txt", []byte(data), 0644)
}


func receive() {
	conn, _ := net.ListenPacket(proto, addr)
	
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		buf := make([]byte, 1024)
		num_of_bytes, source, err := conn.ReadFrom(buf)
		//HER HAR VI IF, IKKE DREP OSS PLZ <3
		if err != nil {
			fmt.Println("NO DATA RECEIVED FOR 5 SECONDS")
		}else {
			fmt.Println(string(buf[:num_of_bytes]))
			conn.WriteTo(buf,source)
		}


	}
}



