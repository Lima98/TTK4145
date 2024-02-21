package main

import (
	"Elevator_project/driver-go/elevio"
	"Elevator_project/fsm"
)

func main(){

    numFloors := 4

    elevio.Init("localhost:15657", numFloors)
    
    var d elevio.MotorDirection = elevio.MD_Stop
    elevio.SetMotorDirection(d)
    
    go fsm.Statemachine()
    
    select {}
      
}
