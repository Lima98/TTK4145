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
	//Forsøk på å kjøre 3

	var id string
	numFloors := 4
	id = os.Args[1]
	fmt.Println("id: ", id)

	switch id {
	case "0":
		elevio.Init("localhost:15657", numFloors)
		go autorestart.ProcessPair(proto, "localhost:20022", "localhost:30022", id)
	case "1":
		elevio.Init("localhost:15658", numFloors)
		go autorestart.ProcessPair(proto, "localhost:20023", "localhost:30023", id)
	case "2":
		elevio.Init("localhost:15659", numFloors)
		go autorestart.ProcessPair(proto, "localhost:20024", "localhost:30024", id)
	}
	var d elevio.MotorDirection = elevio.MD_Stop
	elevio.SetMotorDirection(d)
	//Kjøre en heis
	// numFloors := 4
	// elevio.Init("localhost:15657",numFloors)
	// var d elevio.MotorDirection = elevio.MD_Stop
	// elevio.SetMotorDirection(d)
	// go autorestart.ProcessPair(proto, addrFsmPp, addrPpBackup, "0")

	select {}

}
