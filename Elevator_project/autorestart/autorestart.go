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

const backupFilePath = "./autorestart/cab_orders.txt"


type ProgramType int
const (
	Primary 	= 0
	Backup	= 1
)



func ProcessPair(proto string, addrFsmPp string, addrPpBackup string, mode string, id string) {

	var programtype = Backup

	data, _ := os.ReadFile(backupFilePath)
	var pid string

	// BACKUP MÅ OGSÅ FØRSØKE MORD PÅ PRIMARY VED OVERTAKELSE/KUPP

	//data := []byte{0, 0, 0, 0}
	fmt.Print("FROM THE FILE WE READ: ")
	fmt.Println(data)

	for {
		switch programtype {
		case Primary:
			conn, err := net.Dial(proto, addrPpBackup)
			conn.Write([]byte(strconv.Itoa(os.Getpid())))
			if err == nil {
				conn.Close()
			}
		case Backup:
			if checkMaster(proto, addrPpBackup) {
				conn, _ := net.ListenPacket(proto, addrPpBackup)
				conn.SetReadDeadline(time.Now().Add(5 * time.Second))
				buf := make([]byte, 1024)
				num_of_bytes, _, _ := conn.ReadFrom(buf)
				pid = string(buf[:num_of_bytes])
				fmt.Println("My master is " + pid)
				conn.Close()
			} else {
				exec.Command("gnome-terminal", "--", "kill", "-TERM", pid).Run() //opens a new window so might be messy
				programtype = Primary
				data, _ := os.ReadFile(backupFilePath)
				if data == nil {
					data = []byte{0, 0, 0, 0}
				}
				go fsm.Statemachine(proto, addrFsmPp, data, id)
				exec.Command("gnome-terminal", "--", "go", "run", "./main.go", mode, id).Run()
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
