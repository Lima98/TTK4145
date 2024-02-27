package autorestart

import (
	"Elevator_project/fsm"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// Spør studass om hvordan vi kan definere disse en gang når de brukes på tvers av filer

const cabOrders = "./cab_orders.txt"

func ProcessPair(proto string, addrFsmPp string,addrPpBackup string) {

	var programtype = 1 //0 is primary, 1 is backup

	data, _ := os.ReadFile(cabOrders)

	for {
		switch programtype {
		case 0:
			data = receive(proto, addrFsmPp)
			send(proto, addrPpBackup)

			fmt.Println("I AM THE MASTER NOW")
			writeToFile(data)

		case 1:
			if checkMaster(proto, addrPpBackup) {
				fmt.Println("I received from master")
			}else {
				go fsm.Statemachine()
				createBackup()

				programtype = 0
				//data, _ = os.ReadFile(cabOrders)
			}
		}
	}
}

func checkMaster(proto string, addr string) bool{
	conn, _ := net.ListenPacket(proto, addr)

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	buf := make([]byte, 1024)
	_, _, err := conn.ReadFrom(buf)

	if err != nil {
		conn.Close()
		return false
	}else {
		conn.Close()
		return true
	}	
}

func send(proto string, addr string){
	conn, _ := net.Dial(proto,addr)                                                                                  
	conn.Write([]byte("I am alive"))                                                                               
}

func readFromFile() int {
	data, _ := os.ReadFile(cabOrders)
	counter, _ := strconv.Atoi(string(data))
	return counter
}

func writeToFile(data []byte){
	os.WriteFile("./cab_orders.txt",  data, 0644)
}

func receive(proto string, addr string) []byte{
	conn, _ := net.ListenPacket(proto, addr)

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 1024)
	num_of_bytes, _, _ := conn.ReadFrom(buf)
	fmt.Println(string(buf[:num_of_bytes]))
	conn.Close()
	return buf

}

func createBackup() {
	exec.Command("gnome-terminal", "--", "go", "run", "./main.go").Run()
}


