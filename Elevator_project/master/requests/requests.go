package requests

import (
	elev "Elevator_project/master/elevator"
	elevio "Elevator_project/master/driver-go/elevio"
)


type DirnBehaviourPair struct {
	Dir			elev.Dirn
	Behaviour 	elev.ElevatorBehaviour
}

func requests_above(e elev.Elevator) bool{
	for f := e.Floor+1; f < elev.N_FLOORS; f++ {
		for btn := 0; btn < elev.N_BUTTONS; btn++ {
			if(e.Requests[f][btn]){
				return true
			}
		}
	}
	return false
}

func requests_below(e elev.Elevator) bool{
	for f := 0; f < elev.N_FLOORS; f++ {
		for btn := 0; btn < elev.N_BUTTONS; btn++ {
			if(e.Requests[f][btn]){
				return true
			}
		}
	}
	return false
}

func requests_here(e elev.Elevator) bool{
	for btn := 0; btn < elev.N_BUTTONS; btn++ {
		if(e.Requests[e.Floor][btn]){
			return true
		}
	}
	return false
}

func requests_chooseDirection(e elev.Elevator) DirnBehaviourPair {
		switch e.Dir {
		case elev.D_Up:
			switch {
			case requests_above(e):
				return DirnBehaviourPair{elev.D_Up, elev.EB_Moving}
			case requests_here(e):
				return DirnBehaviourPair{elev.D_Down, elev.EB_DoorOpen}
			case requests_below(e):
				return DirnBehaviourPair{elev.D_Down, elev.EB_Moving}
			default:
				return DirnBehaviourPair{elev.D_Stop, elev.EB_Idle}
			}
		case elev.D_Down:
			switch {
			case requests_below(e):
				return DirnBehaviourPair{elev.D_Down, elev.EB_Moving}
			case requests_here(e):
				return DirnBehaviourPair{elev.D_Up, elev.EB_DoorOpen}
			case requests_above(e):
				return DirnBehaviourPair{elev.D_Up, elev.EB_Moving}
			default:
				return DirnBehaviourPair{elev.D_Stop, elev.EB_Idle}
			}
		case elev.D_Stop:
			switch {
			case requests_here(e):
				return DirnBehaviourPair{elev.D_Stop, elev.EB_DoorOpen}
			case requests_above(e):
				return DirnBehaviourPair{elev.D_Up, elev.EB_Moving}
			case requests_below(e):
				return DirnBehaviourPair{elev.D_Down, elev.EB_Moving}
			default:
				return DirnBehaviourPair{elev.D_Stop, elev.EB_Idle}
			}
		default:
			return DirnBehaviourPair{elev.D_Stop, elev.EB_Idle}
		}
	}

func requests_shouldStop(e elev.Elevator) bool{
		switch(e.Dir){
		case elev.D_Down:
			return bool(e.Requests[e.Floor][elevio.BT_HallDown] ||
						e.Requests[e.Floor][elevio.BT_Cab]||
						!requests_below(e))
		case elev.D_Up:
			return e.Requests[e.Floor][elevio.BT_HallUp]   ||
				   e.Requests[e.Floor][elevio.BT_Cab]      ||
				   !requests_above(e)
		case elev.D_Stop:
			return false
		default:
			return true
		}
	}	

	