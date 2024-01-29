package main

import (
	"fmt"
)

const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

type Dirn int

const (
	D_Up Dirn = iota
	D_Down
	D_Stop
)

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

type Button int

const (
	B_HallUp Button = iota
	B_HallDown
	B_Cab
)

type ClearRequestVariant int

const (
	CV_All ClearRequestVariant = iota
	CV_InDirn
)

type Elevator struct {
	floor     int
	dirn      Dirn
	requests  [N_FLOORS][N_BUTTONS]int
	behaviour ElevatorBehaviour
	config    struct {
		clearRequestVariant ClearRequestVariant
		doorOpenDurationS   float64
	}
}

func elevatorPrint(e Elevator) {
	fmt.Println("  +--------------------+")
	fmt.Printf(
		"  |floor = %-2d          |\n"+
			"  |dirn  = %-12.12s|\n"+
			"  |behav = %-12.12s|\n",
		e.floor,
		dirnToString(e.dirn),
		ebToString(e.behaviour),
	)
	fmt.Println("  +--------------------+")
	fmt.Println("  |  | up  | dn  | cab |")
	for f := N_FLOORS - 1; f >= 0; f-- {
		fmt.Printf("  | %d", f)
		for btn := 0; btn < N_BUTTONS; btn++ {
			if (f == N_FLOORS-1 && btn == int(B_HallUp)) ||
				(f == 0 && btn == int(B_HallDown)) {
				fmt.Print("|     ")
			} else {
				fmt.Printf(e.requests[f][btn]*3, "|  #  ", "|  -  ")
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +--------------------+")
}

func ebToString(eb ElevatorBehaviour) string {
	switch eb {
	case EB_Idle:
		return "EB_Idle"
	case EB_DoorOpen:
		return "EB_DoorOpen"
	case EB_Moving:
		return "EB_Moving"
	default:
		return "EB_UNDEFINED"
	}
}

func dirnToString(dirn Dirn) string {
	switch dirn {
	case D_Up:
		return "D_Up"
	case D_Down:
		return "D_Down"
	case D_Stop:
		return "D_Stop"
	default:
		return "D_UNDEFINED"
	}
}

func elevatorUninitialized() Elevator {
	return Elevator{
		floor:     -1,
		dirn:      D_Stop,
		behaviour: EB_Idle,
		config: struct {
			clearRequestVariant ClearRequestVariant
			doorOpenDurationS   float64
		}{
			clearRequestVariant: CV_All,
			doorOpenDurationS:   3.0,
		},
	}
}
