package requests

import (
	elevio "Elevator_project/driver-go/elevio"
	elev "Elevator_project/elevator"
	"fmt"
)


type DirBehaviourPair struct {
	Dir			elevio.MotorDirection
	Behaviour 	elev.ElevatorBehaviour
}

func Requests_above(e elev.Elevator) bool{
	for f := e.Floor+1; f < elev.N_FLOORS; f++ {
		for btn := 0; btn < elev.N_BUTTONS; btn++ {
			if(e.Requests[f][btn]){
				return true
			}
		}
	}
	return false
}

func Requests_below(e elev.Elevator) bool{
	for f := 0; f < elev.N_FLOORS; f++ {
		for btn := 0; btn < elev.N_BUTTONS; btn++ {
			if(e.Requests[f][btn]){
				return true
			}
		}
	}
	return false
}

func Requests_here(e elev.Elevator) bool{
	for btn := 0; btn < elev.N_BUTTONS; btn++ {
		if(e.Requests[e.Floor][btn]){
			return true
		}
	}
	return false
}

func Requests_chooseDirection(e elev.Elevator) DirBehaviourPair{
		switch e.Dir {
		case elevio.MD_Up:
			switch {
			case Requests_above(e):
				return DirBehaviourPair{elevio.MD_Up, elev.EB_Moving}
			case Requests_here(e):
				return DirBehaviourPair{elevio.MD_Down, elev.EB_DoorOpen}
			case Requests_below(e):
				return DirBehaviourPair{elevio.MD_Down, elev.EB_Moving}
			default:
				return DirBehaviourPair{elevio.MD_Stop, elev.EB_Idle}
			}
		case elevio.MD_Down:
			switch {
			case Requests_below(e):
				return DirBehaviourPair{elevio.MD_Down, elev.EB_Moving}
			case Requests_here(e):
				return DirBehaviourPair{elevio.MD_Up, elev.EB_DoorOpen}
			case Requests_above(e):
				return DirBehaviourPair{elevio.MD_Up, elev.EB_Moving}
			default:
				return DirBehaviourPair{elevio.MD_Stop, elev.EB_Idle}
			}
		case elevio.MD_Stop:
			switch {
			case Requests_here(e):
				return DirBehaviourPair{elevio.MD_Stop, elev.EB_DoorOpen}
			case Requests_above(e):
				return DirBehaviourPair{elevio.MD_Up, elev.EB_Moving}
			case Requests_below(e):
				return DirBehaviourPair{elevio.MD_Down, elev.EB_Moving}
			default:
				return DirBehaviourPair{elevio.MD_Stop, elev.EB_Idle}
			}
		default:
			return DirBehaviourPair{elevio.MD_Stop, elev.EB_Idle}
		}
	}

func Requests_shouldStop(e elev.Elevator) bool{
		switch(e.Dir){
		case elevio.MD_Down:
			return bool(e.Requests[e.Floor][elevio.BT_HallDown] ||
						e.Requests[e.Floor][elevio.BT_Cab] 		||
						!Requests_below(e))
		case elevio.MD_Up:
			return e.Requests[e.Floor][elevio.BT_HallUp]   ||
				   e.Requests[e.Floor][elevio.BT_Cab]      ||
				   !Requests_above(e)
		case elevio.MD_Stop:
			return false
		default:
			return true
		}
	}	

func Requests_shouldClearImmediately(e elev.Elevator, btn_floor int, btn_type elevio.ButtonType) bool{
	return bool(e.Floor == btn_floor && (
		(e.Dir == elevio.MD_Up   && btn_type == elevio.BT_HallUp)    ||
		(e.Dir == elevio.MD_Down && btn_type == elevio.BT_HallDown)  ||
		e.Dir == elevio.MD_Stop || btn_type == elevio.BT_Cab))  
}


func Requests_clearAtCurrentFloor(e elev.Elevator) elev.Elevator{  
        e.Requests[e.Floor][elevio.BT_Cab] = false
        switch(e.Dir){
        case elevio.MD_Up:
            if(!Requests_above(e) && !e.Requests[e.Floor][elevio.BT_HallUp]){
                e.Requests[e.Floor][elevio.BT_HallDown] = false
            }
            e.Requests[e.Floor][elevio.BT_HallUp] = false
            
        case elevio.MD_Down:
            if(!Requests_below(e) && !e.Requests[e.Floor][elevio.BT_HallDown]){
                e.Requests[e.Floor][elevio.BT_HallUp] = false
            }
            e.Requests[e.Floor][elevio.BT_HallDown] = false
            
        case elevio.MD_Stop:
        default:
            e.Requests[e.Floor][elevio.BT_HallUp] = false
            e.Requests[e.Floor][elevio.BT_HallDown] = false
        }
		return e
}
	

func PrintRequests(e elev.Elevator) {
	fmt.Println(e.Requests)
}