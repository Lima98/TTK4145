package main

func requestsChooseDirection(e Elevator) DirnBehaviourPair {
	switch e.dirn {
	case D_Up:
		return DirnBehaviourPair{
			D_Up, EB_Moving,
		}
	case D_Down:
		return DirnBehaviourPair{
			D_Down, EB_Moving,
		}
	case D_Stop:
		return DirnBehaviourPair{
			D_Stop, EB_Idle,
		}
	default:
		return DirnBehaviourPair{
			D_Stop, EB_Idle,
		}
	}
}

func requestsShouldStop(e Elevator) bool {
	switch e.dirn {
	case D_Down:
		return e.requests[e.floor][B_HallDown] ||
			e.requests[e.floor][B_Cab] ||
			!requestsBelow(e)
	case D_Up:
		return e.requests[e.floor][B_HallUp] ||
			e.requests[e.floor][B_Cab] ||
			!requestsAbove(e)
	case D_Stop:
		return true
	default:
		return true
	}
}

func requestsShouldClearImmediately(e Elevator, btnFloor int, btnType Button) bool {
	switch e.config.clearRequestVariant {
	case CV_All:
		return e.floor == btnFloor
	case CV_InDirn:
		return e.floor == btnFloor &&
			((e.dirn == D_Up && btnType == B_HallUp) ||
				(e.dirn == D_Down && btnType == B_HallDown) ||
				e.dirn == D_Stop ||
				btnType == B_Cab)
	default:
		return false
	}
}

func requestsClearAtCurrentFloor(e Elevator) Elevator {
	switch e.config.clearRequestVariant {
	case CV_All:
		for btn := 0; btn < N_BUTTONS; btn++ {
			e.requests[e.floor][btn] = 0
		}
	case CV_InDirn:
		e.requests[e.floor][B_Cab] = 0
		switch e.dirn {
		case D_Up:
			if !requestsAbove(e) && !e.requests[e.floor][B_HallUp] {
				e.requests[e.floor][B_HallDown] = 0
			}
			e.requests[e.floor][B_HallUp] = 0
		case D_Down:
			if !requestsBelow(e) && !e.requests[e.floor][B_HallDown] {
				e.requests[e.floor][B_HallUp] = 0
			}
			e.requests[e.floor][B_HallDown] = 0
		case D_Stop:
			e.requests[e.floor][B_HallUp] = 0
			e.requests[e.floor][B_HallDown] = 0
		}
	default:
		// handle other cases
	}
	return e
}

func requestsAbove(e Elevator) bool {
	for f := e.floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.requests[f][btn] != 0 {
				return true
			}
		}
	}
	return false
}

func requestsBelow(e Elevator) bool {
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.requests[f][btn] != 0 {
				return true
			}
		}
	}
	return false
}

func requestsHere(e Elevator) bool {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if e.requests[e.floor][btn] != 0 {
			return true
		}
	}
	return false
}
