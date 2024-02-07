package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const proto, addr = "udp", "localhost:20022"
const dataFile = "./data.txt"
var programtype = 1 //0 is primary, 1 is backup


func main() {
	data, _ := os.ReadFile(dataFile)
	counter, _ := strconv.Atoi(string(data))

	for {
		switch programtype {
		case 0:
			counter = readFromFile()
			send()
			time.Sleep(1 * time.Second)
			counter++
			fmt.Println("I AM THE MASTER NOW")
			fmt.Println(counter)
			writeToFile(strconv.FormatInt(int64(counter),10))
		case 1:
			receive()
		}
	}
}

func send(){
	conn, _ := net.Dial(proto,addr)                                                                                  
	conn.Write([]byte("I am alive"))                                                                               
}

func readFromFile() int {
	data, _ := os.ReadFile(dataFile)
	counter, _ := strconv.Atoi(string(data))
	return counter
}

func writeToFile(data string){
	os.WriteFile("./data.txt", []byte(data), 0644)
}

func receive() {
	conn, _ := net.ListenPacket(proto, addr)
	
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		buf := make([]byte, 1024)
		num_of_bytes, source, err := conn.ReadFrom(buf)
		//HER HAR VI IF, IKKE DREP OSS PLZ <3
		if err != nil {
			conn.Close()
			fmt.Println("NO MASTER FOUND")
			programtype = 0
			fmt.Println("Going master and making backup")
			createBackup()
		}else {
			fmt.Println(string(buf[:num_of_bytes]))
			conn.WriteTo(buf,source)
			fmt.Println("-received from master")
			programtype = 1
			conn.Close()
		}

}

func createBackup() {
	exec.Command("gnome-terminal", "--", "go", "run", "B.go").Run()
}
