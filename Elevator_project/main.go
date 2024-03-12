package main

import (
	"Elevator_project/autorestart"
	"Elevator_project/driver-go/elevio"
	"fmt"
	"os"
)

const proto, addrFsmPp = "udp", "localhost:20022"
const addrPpBackup = "localhost:30022"

func main() {

	var id string
	var mode string
	numFloors := 4

	mode = os.Args[1]
	id = os.Args[2]

	fmt.Println("id: ", id)

	if mode == "sim" {
		switch id {
		case "0":
			elevio.Init("localhost:15657", numFloors)
			go autorestart.ProcessPair(proto, "localhost:20022", "localhost:30022", mode, id)
		case "1":
			elevio.Init("localhost:15658", numFloors)
			go autorestart.ProcessPair(proto, "localhost:20023", "localhost:30023", mode, id)
		case "2":
			elevio.Init("localhost:15659", numFloors)
			go autorestart.ProcessPair(proto, "localhost:20024", "localhost:30024", mode, id)
		}
	}
	if mode == "physical" {
		elevio.Init("localhost:15657", numFloors)
		go autorestart.ProcessPair(proto, "localhost:20022", "localhost:30022", mode, id)
	}

	var d elevio.MotorDirection = elevio.MD_Stop
	elevio.SetMotorDirection(d)

	select {}

}
