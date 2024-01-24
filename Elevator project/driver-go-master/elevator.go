package main

import (
	"Driver-go/elevio"
)

type ElevatorBehaviour int

const (
	EB_Idle		ElevatorBehaviour 	= 0
    EB_DoorOpen 				 	= -1
    EB_Moving						= 1
)

type Elevator struct {
	floor int
    dirn elevio.MotorDirection
    //int                     requests[N_FLOORS][N_BUTTONS];
    behaviour ElevatorBehaviour
	doorOpenDuration_s float32

}

func elevator_uninitialized() Elevator{
    return Elevator{
        floor: -1,
        dirn: elevio.MD_Stop,
        behaviour: EB_Idle,
        doorOpenDuration_s: 3.0,
    }
}