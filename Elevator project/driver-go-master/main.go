package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Started!")

	var inputPollRateMs = 25
	conLoad("elevator.con",
		conVal(&inputPollRateMs, "%d"),
	)

	var input ElevInputDevice // Replace with your actual implementation

	elevator := &Elevator{
		timer: newTimer(),
	}

	if input.floorSensor() == -1 {
		elevator.fsmOnInitBetweenFloors()
	}

	for {
		// Request button
		for f := 0; f < N_FLOORS; f++ {
			for b := 0; b < N_BUTTONS; b++ {
				v := input.requestButton(f, b)
				if v != 0 && v != elevator.prevRequests[f][b] {
					elevator.fsmOnRequestButtonPress(f, b)
				}
				elevator.prevRequests[f][b] = v
			}
		}

		// Floor sensor
		f := input.floorSensor()
		if f != -1 && f != elevator.prevFloor {
			elevator.fsmOnFloorArrival(f)
		}
		elevator.prevFloor = f

		// Timer
		if elevator.timer.timedOut() {
			elevator.timer.stop()
			elevator.fsmOnDoorTimeout()
		}

		time.Sleep(time.Duration(inputPollRateMs) * time.Millisecond)
	}
}

func (e *Elevator) fsmOnInitBetweenFloors() {
	// Placeholder for fsm_onInitBetweenFloors function in C
}

func (e *Elevator) fsmOnRequestButtonPress(floor, button int) {
	// Placeholder for fsm_onRequestButtonPress function in C
}

func (e *Elevator) fsmOnFloorArrival(floor int) {
	// Placeholder for fsm_onFloorArrival function in C
}

func (e *Elevator) fsmOnDoorTimeout() {
	// Placeholder for fsm_onDoorTimeout function in C
}
