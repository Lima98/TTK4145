package elevator

var N_FLOORS = 4
var N_BUTTONS = 3

type ElevatorBehaviour int
const (
	EB_Idle	ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

type Dirn int
const (
	D_Down 		Dirn = -1
	D_Stop 		Dirn = 0
	D_Up		Dirn = 1
)

type Request struct {
	assigned_to	int
	completed	int
}

type Elevator struct {
	floor	int
	dir		Dirn
	request Request
	behaviour ElevatorBehaviour
}

