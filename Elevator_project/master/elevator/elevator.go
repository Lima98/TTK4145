package elevator

const N_FLOORS = 4
const N_BUTTONS = 3

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
	Requested	int
	Assigned_To	int
	Completed	int
}

type Elevator struct {
	Floor	int
	Dir		Dirn
	Requests [N_FLOORS][N_BUTTONS]bool
	Behaviour ElevatorBehaviour
}

