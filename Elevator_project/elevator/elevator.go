package elevator

import elevio "Elevator_project/driver-go/elevio"

const N_FLOORS = 4
const N_BUTTONS = 3

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
}

// func NewElevator(elev Elevator) *Elevator {
// 	return &Elevator{Floor: -1, Dir: MD_Stop, Requests: [N_FLOORS][N_BUTTONS], Behaviour: EB_Idle}
// }