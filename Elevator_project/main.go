package main

import (
	"Elevator_project/autorestart"
	"Elevator_project/driver-go/elevio"
)

const proto, addrFsmPp = "udp", "localhost:20022"
const addrPpBackup = "localhost:30022"



func main(){

    numFloors := 4

    elevio.Init("localhost:15657", numFloors)
    
    var d elevio.MotorDirection = elevio.MD_Stop
    elevio.SetMotorDirection(d)
    
    go autorestart.ProcessPair(proto, addrFsmPp,addrPpBackup)
    
    select {}

}     
