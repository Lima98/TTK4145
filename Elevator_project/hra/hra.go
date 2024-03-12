package hra

import (
	elevio "Elevator_project/driver-go/elevio"
	elev "Elevator_project/elevator"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAElevState struct {
	Behavior    string              `json:"behaviour"`
	Floor       int                 `json:"floor"`
	Direction   string              `json:"direction"`
	CabRequests [elev.N_FLOORS]bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [elev.N_FLOORS][2]bool  `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func HallRequestAssigner(orders [elev.N_FLOORS][elev.N_BUTTONS - 1]int, Elevators [elev.N_ELEVATORS]elev.Elevator) map[string][][2]bool {

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	HallRequestsTemp := [elev.N_FLOORS][2]bool{}

	for i := 0; i < elev.N_FLOORS; i++ {
		for j := 0; j < elev.N_BUTTONS-1; j++ {
			if orders[i][j] < 2 {
				HallRequestsTemp[i][j] = true
			} else {
				HallRequestsTemp[i][j] = false
			}
		}
	}

	CabRequestsTemp1 := [elev.N_FLOORS]bool{}

	for i := 0; i < elev.N_FLOORS; i++ {
		CabRequestsTemp1[i] = Elevators[0].Requests[i][2]
	}

	CabRequestsTemp2 := [elev.N_FLOORS]bool{}

	for i := 0; i < elev.N_FLOORS; i++ {
		CabRequestsTemp2[i] = Elevators[1].Requests[i][2]
	}

	CabRequestsTemp3 := [elev.N_FLOORS]bool{}

	for i := 0; i < elev.N_FLOORS; i++ {
		CabRequestsTemp3[i] = Elevators[2].Requests[i][2]
	}

	behaviorToString := make(map[elev.ElevatorBehaviour]string)
	behaviorToString[elev.EB_DoorOpen] = "idle"
	behaviorToString[elev.EB_Idle] = "idle"
	behaviorToString[elev.EB_Moving] = "moving"

	directionToString := make(map[elevio.MotorDirection]string)
	directionToString[elevio.MD_Up] = "up"
	directionToString[elevio.MD_Down] = "down"
	directionToString[elevio.MD_Stop] = "stop"

	StatesTemp := map[string]HRAElevState{
		"one": HRAElevState{
			Behavior:    behaviorToString[Elevators[0].Behaviour],
			Floor:       Elevators[0].Floor,
			Direction:   directionToString[Elevators[0].Dir],
			CabRequests: CabRequestsTemp1,
		},
		"two": HRAElevState{
			Behavior:    behaviorToString[Elevators[1].Behaviour],
			Floor:       Elevators[1].Floor,
			Direction:   directionToString[Elevators[1].Dir],
			CabRequests: CabRequestsTemp1,
		},
		"three": HRAElevState{
			Behavior:    behaviorToString[Elevators[2].Behaviour],
			Floor:       Elevators[2].Floor,
			Direction:   directionToString[Elevators[2].Dir],
			CabRequests: CabRequestsTemp1,
		},
	}

	input := HRAInput{
		HallRequests: HallRequestsTemp,
		States:       StatesTemp,
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
	}

	ret, err := exec.Command(hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
	}

	fmt.Printf("output: \n")
	for k, v := range *output {
		fmt.Printf("%6v :  %+v\n", k, v)
	}
	return *output
}
