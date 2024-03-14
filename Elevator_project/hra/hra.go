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

func HallRequestAssigner(orders [elev.N_FLOORS][elev.N_BUTTONS - 1]int, Elevators map[string]elev.Elevator, peers []string) map[string][][2]bool {

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	behaviorToString := make(map[elev.ElevatorBehaviour]string)
	behaviorToString[elev.EB_DoorOpen] = "idle"
	behaviorToString[elev.EB_Idle] = "idle"
	behaviorToString[elev.EB_Moving] = "moving"

	directionToString := make(map[elevio.MotorDirection]string)
	directionToString[elevio.MD_Up] = "up"
	directionToString[elevio.MD_Down] = "down"
	directionToString[elevio.MD_Stop] = "stop"

	HallRequestsTemp := [elev.N_FLOORS][2]bool{}

	for i := 0; i < elev.N_FLOORS; i++ {
		for j := 0; j < elev.N_BUTTONS-1; j++ {
			if orders[i][j] < elev.Completed {
				HallRequestsTemp[i][j] = true
			} else {
				HallRequestsTemp[i][j] = false
			}
		}
	}

	NumPeers := len(peers)
	CabRequestsTemp := make([][elev.N_FLOORS]bool, NumPeers)

	for i := 0; i < len(peers); i++ {
		for j := 0; j < elev.N_FLOORS; j++ {
			CabRequestsTemp[i][j] = Elevators[peers[i]].Requests[j][2]
		}
	}

	StatesTemp := map[string]HRAElevState{}

	for i := 0; i < len(peers); i++ {
		if Elevators[peers[i]].Obstructed {
			StatesTemp[peers[i]] = HRAElevState{}
		}else{
			StatesTemp[peers[i]] = HRAElevState{
				Behavior:    behaviorToString[Elevators[peers[i]].Behaviour],
				Floor:       Elevators[peers[i]].Floor,
				Direction:   directionToString[Elevators[peers[i]].Dir],
				CabRequests: CabRequestsTemp[i],
			}
		}
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
