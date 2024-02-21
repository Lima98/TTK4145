package elevator

import (
	elevio "Elevator_project/driver-go/elevio"
	"time"
)

const N_FLOORS = 4
const N_BUTTONS = 3
const OPEN_DOOR_TIME = 3 *time.Second


type ElevatorBehaviour int
const (
	EB_Idle	ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

type Request struct {
	Requested	int
	Assigned_To	int
	Completed	int
}

type Elevator struct {
	Floor	int
	Dir		elevio.MotorDirection
	Requests [N_FLOORS][N_BUTTONS]bool
	Behaviour ElevatorBehaviour
	Obstructed bool
}

func PrintBehaviour(e Elevator) {
	switch e.Behaviour {
	case EB_Idle:
		println("IDLE")
	case EB_DoorOpen:
		println("DOOR OPEN")
	case EB_Moving:
	println("MOVING")
	}
}

