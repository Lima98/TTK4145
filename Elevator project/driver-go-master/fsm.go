package main

import (
	"fmt"
)

type DirnBehaviourPair struct {
	dirn      Dirn
	behaviour ElevatorBehaviour
}

func (t *Timer) start(duration float64) {
	t.endTime = getWallTime() + duration
	t.active = true
	t.duration = duration
}

func (t *Timer) stop() {
	t.active = false
}

func (t *Timer) timedOut() bool {
	return t.active && getWallTime() > t.endTime
}

func getWallTime() float64 {
	// Placeholder for get_wall_time function in C
	// You can implement this function using the time package in Go
	return 0
}

func (e *Elevator) print() {
	// Placeholder for elevator_print function in C
	// You can implement this function to print the elevator state in Go
}

func fsmOnInitBetweenFloors(outputDevice ElevOutputDevice) Elevator {
	outputDevice.motorDirection(D_Down)
	return Elevator{dirn: D_Down, behaviour: EB_Moving}
}

func fsmOnRequestButtonPress(outputDevice ElevOutputDevice, btnFloor int, btnType Button, elevator Elevator) Elevator {
	fmt.Printf("\n\n%s(%d, %s)\n", "fsm_onRequestButtonPress", btnFloor, buttonToString(btnType))
	elevator.print()

	switch elevator.behaviour {
	case EB_DoorOpen:
		if requestsShouldClearImmediately(elevator, btnFloor, btnType) {
			newTimer().start(elevator.config.doorOpenDurationS)
		} else {
			elevator.requests[btnFloor][btnType] = 1
		}
	case EB_Moving:
		elevator.requests[btnFloor][btnType] = 1
	case EB_Idle:
		elevator.requests[btnFloor][btnType] = 1
		pair := requestsChooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour
		switch pair.behaviour {
		case EB_DoorOpen:
			outputDevice.doorLight(1)
			newTimer().start(elevator.config.doorOpenDurationS)
			elevator = requestsClearAtCurrentFloor(elevator)
		case EB_Moving:
			outputDevice.motorDirection(elevator.dirn)
		case EB_Idle:
		}
	}
	setAllLights(outputDevice, elevator)

	fmt.Printf("\nNew state:\n")
	elevator.print()

	return elevator
}

func fsmOnFloorArrival(outputDevice ElevOutputDevice, newFloor int, elevator Elevator) Elevator {
	fmt.Printf("\n\n%s(%d)\n", "fsm_onFloorArrival", newFloor)
	elevator.print()

	elevator.floor = newFloor

	outputDevice.floorIndicator(elevator.floor)

	switch elevator.behaviour {
	case EB_Moving:
		if requestsShouldStop(elevator) {
			outputDevice.motorDirection(D_Stop)
			outputDevice.doorLight(1)
			elevator = requestsClearAtCurrentFloor(elevator)
			newTimer().start(elevator.config.doorOpenDurationS)
			setAllLights(outputDevice, elevator)
			elevator.behaviour = EB_DoorOpen
		}
	default:
	}

	fmt.Printf("\nNew state:\n")
	elevator.print()

	return elevator
}

func fsmOnDoorTimeout(outputDevice ElevOutputDevice, elevator Elevator) Elevator {
	fmt.Printf("\n\n%s()\n", "fsm_onDoorTimeout")
	elevator.print()

	switch elevator.behaviour {
	case EB_DoorOpen:
		pair := requestsChooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour

		switch elevator.behaviour {
		case EB_DoorOpen:
			newTimer().start(elevator.config.doorOpenDurationS)
			elevator = requestsClearAtCurrentFloor(elevator)
			setAllLights(outputDevice, elevator)
		case EB_Moving, EB_Idle:
			outputDevice.doorLight(0)
			outputDevice.motorDirection(elevator.dirn)
		}
	default:
	}

	fmt.Printf("\nNew state:\n")
	elevator.print()

	return elevator
}

func main() {
	fmt.Println("Started!")

	outputDevice := elevioGetOutputDevice()
	timer := newTimer()
	elevator := elevatorUninitialized()

	conLoad("elevator.con",
		conVal(&elevator.config.doorOpenDurationS, "%lf"),
		conEnum(&elevator.config.clearRequestVariant,
			conMatch(CV_All),
			conMatch(CV_InDirn),
		),
	)

	outputDevice = elevioGetOutputDevice()

	elevator = fsmOnInitBetweenFloors(outputDevice)

	for {
		// Placeholder for the main loop logic
	}
}

func setAllLights(outputDevice ElevOutputDevice, elevator Elevator) {
	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			outputDevice.requestButtonLight(floor, btn, elevator.requests[floor][btn])
		}
	}
}

func buttonToString(btnType Button) string {
	// Placeholder for buttonToString function in C
	// You can implement this function to convert Button type to string in Go
	return ""
}
