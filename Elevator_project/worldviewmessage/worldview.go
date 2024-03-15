package worldviewmessage

import (
	"Elevator_project/driver-go/elevio"
	elev "Elevator_project/elevator"
)

type Order struct {
	Order             elev.OrderState
	ElevatorsThatKnow map[string]bool // med id
}

type WorldViewMsg struct {
	Orders        [elev.N_FLOORS][elev.N_BUTTONS - 1]Order
	ID            string
	Fault         bool
	ElevatorState elev.Elevator
}

type WorldView struct {
	Orders    [elev.N_FLOORS][elev.N_BUTTONS - 1]Order
	Elevators map[string]elev.Elevator
}


func Orders_clearAtCurrentFloor(wv WorldView, e elev.Elevator) WorldView {
	if !e.Requests[e.Floor][elevio.BT_HallDown] {
		wv.Orders[e.Floor][elevio.BT_HallDown].Order = elev.Completed
	}
	if !e.Requests[e.Floor][elevio.BT_HallUp] {
		wv.Orders[e.Floor][elevio.BT_HallUp].Order = elev.Completed
	}
	return wv
} 
