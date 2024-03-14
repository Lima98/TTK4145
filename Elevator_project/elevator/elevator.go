package elevator

import (
	elevio "Elevator_project/driver-go/elevio"
	"time"
)

const N_FLOORS = 4
const N_BUTTONS = 3
const N_ELEVATORS = 3
const OPEN_DOOR_TIME = 3 * time.Second
const FAULT_TIMEOUT = 6 * time.Second

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

type OrderState int

const (
	Unknown    = -1
	Completed  = 0
	Unassigned = 1
	Assigned   = 2
)

type Elevator struct {
	Floor      int
	Dir        elevio.MotorDirection
	Requests   [N_FLOORS][N_BUTTONS]bool
	Behaviour  ElevatorBehaviour
	Obstructed bool
	ID         string
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
