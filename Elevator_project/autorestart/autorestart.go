package autorestart

import (
	"Elevator_project/fsm"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

// Spør studass om hvordan vi kan definere disse en gang når de brukes på tvers av filer

const backupFilePath = "./autorestart/cab_orders.txt"

func ProcessPair(proto string, addrFsmPp string, addrPpBackup string) {

	var programtype = 1 //0 is primary, 1 is backup

	data, _ := os.ReadFile(backupFilePath)

	//data := []byte{0, 0, 0, 0}
	fmt.Print("FROM THE FILE WE READ: ")
	fmt.Println(data)

	for {
		switch programtype {
		case 0:
			conn, err := net.Dial(proto, addrPpBackup)
			conn.Write([]byte("I am alive"))
			if err == nil {
				conn.Close()
			}
		case 1:
			if checkMaster(proto, addrPpBackup) {
				fmt.Println("I received from master")
			} else {

				programtype = 0
				data, _ := os.ReadFile(backupFilePath)
				if data == nil {
					data = []byte{0, 0, 0, 0}
				}

				go fsm.Statemachine(proto, addrFsmPp, data)
				exec.Command("gnome-terminal", "--", "go", "run", "./main.go").Run()
			}
		}
	}
}

func checkMaster(proto string, addr string) bool {
	conn, err := net.ListenPacket(proto, addr)

	conn.SetReadDeadline(time.Now().Add(1000 * time.Millisecond))
	buf := make([]byte, 1024)
	_, _, err1 := conn.ReadFrom(buf)

	if err1 != nil {
		if err == nil {
			conn.Close()
		}
		return false
	} else {
		if err == nil {
			conn.Close()
		}
		return true
	}
}
