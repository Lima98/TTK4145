package fsm

import (
	elevio "Elevator_project/driver-go/elevio"
	elev "Elevator_project/elevator"
	hra "Elevator_project/hra"
	"Elevator_project/network"
	"Elevator_project/network/network/peers"
	"Elevator_project/requests"
	wv "Elevator_project/worldviewmessage"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func Statemachine(proto string, addr string, cabOrders []byte, id string, backupFilePath string) {

	buttons := make(chan elevio.ButtonEvent)
	floors := make(chan int)
	obstruction := make(chan bool)
	stop := make(chan bool)

	go elevio.PollButtons(buttons)
	go elevio.PollFloorSensor(floors)
	go elevio.PollObstructionSwitch(obstruction)
	go elevio.PollStopButton(stop)

	worldViewTx := make(chan wv.WorldViewMsg)
	worldViewRx := make(chan wv.WorldViewMsg)
	peerUpdateCh := make(chan peers.PeerUpdate)

	go network.Network(worldViewTx, worldViewRx, peerUpdateCh, id)

	var elevator = elev.Elevator{
		Floor:      1,
		Dir:        elevio.MD_Stop,
		Behaviour:  elev.EB_Idle,
		Obstructed: false,
		ID:         id}

	var worldView = wv.WorldView{}
	worldView.Orders = [elev.N_FLOORS][elev.N_BUTTONS - 1]wv.Order{}
	worldView.Elevators = make(map[string]elev.Elevator)
	worldView.Elevators[id] = elevator

	var peerList []string

	for i := 0; i < elev.N_FLOORS; i++ {
		for j := 0; j < elev.N_BUTTONS-1; j++ {
			worldView.Orders[i][j].Order = elev.Unknown
			fmt.Println(("Set to unknown"))
			worldView.Orders[i][j].ElevatorsThatKnow = make(map[string]bool)
		}
	}

	for i := 0; i < elev.N_FLOORS; i++ {
		if cabOrders[i] == 1 {
			elevator.Requests[i][2] = true
		} else {
			elevator.Requests[i][2] = false
		}
	}

	openDoorTimer := time.NewTimer(1000 * time.Second)
	faultTimer := time.NewTimer(1000 * time.Second)

	select {
	case <-floors:
	default:
		elevio.SetMotorDirection(elevio.MD_Down)
		elevator.Dir = elevio.MD_Down
		elevator.Behaviour = elev.EB_Moving
	}

	for {

		select {
		// NETWORK
		case a := <-peerUpdateCh:
			fmt.Println("PEER UPDATE")
			fmt.Println(a)
			fmt.Println("-")
			peerList = a.Peers

		case a := <-worldViewRx:

			worldView.Elevators[a.ID] = a.ElevatorState

			fmt.Print("WORLDVIEW RECEIVED: ")
			fmt.Println(a.Orders)
			/* if a.ID == id {
				break
			} */
			for i := 0; i < elev.N_FLOORS; i++ {
				for j := 0; j < elev.N_BUTTONS-1; j++ {
					switch worldView.Orders[i][j].Order {
					case elev.Unknown:
						worldView.Orders[i][j].Order = a.Orders[i][j].Order
						worldView.Orders[i][j].ElevatorsThatKnow = a.Orders[i][j].ElevatorsThatKnow
					case elev.Completed:
						switch a.Orders[i][j].Order {
						case elev.Unassigned:
							worldView.Orders[i][j].Order = a.Orders[i][j].Order
						case elev.Assigned:
							worldView.Orders[i][j].Order = a.Orders[i][j].Order
							worldView.Orders[i][j].ElevatorsThatKnow = a.Orders[i][j].ElevatorsThatKnow
						case elev.Completed:
						}
					case elev.Unassigned:
						switch a.Orders[i][j].Order {
						case elev.Unassigned:
						case elev.Assigned:
							worldView.Orders[i][j].Order = a.Orders[i][j].Order
							worldView.Orders[i][j].ElevatorsThatKnow = a.Orders[i][j].ElevatorsThatKnow
						case elev.Completed:
						}
					case elev.Assigned:
						switch a.Orders[i][j].Order {
						case elev.Unassigned:
						case elev.Assigned:
						case elev.Completed:
							for k := 0; k < len(peerList); k++ {
								if !worldView.Orders[i][j].ElevatorsThatKnow[peerList[k]] {
									break
								}
							}
							worldView.Orders[i][j].Order = a.Orders[i][j].Order
							worldView.Orders[i][j].ElevatorsThatKnow = make(map[string]bool)
							fmt.Println(worldView.Orders[i][j].ElevatorsThatKnow)

						}
					}
				}
			}

			fmt.Println("UPDATED WORLVIEW: ")
			fmt.Println(worldView.Orders)
			hallAssignments := hra.HallRequestAssigner(worldView.Orders, worldView.Elevators, peerList)

			fmt.Println("HALL ASSIGNMENTS:")
			fmt.Println(hallAssignments)

			for i := 0; i < elev.N_FLOORS; i++ {
				for j := 0; j < elev.N_BUTTONS-1; j++ {
					elevator.Requests[i][j] = hallAssignments[elevator.ID][i][j]
					worldView.Orders[i][j].ElevatorsThatKnow[elevator.ID] = true
					if hallAssignments[elevator.ID][i][j] {
						worldView.Orders[i][j].Order = elev.Assigned
					}
				}
			}
			fmt.Println("Elevator " + id + " has requests")
			fmt.Print(elevator.Requests)
			SetAllLights(worldView, elevator)

		// ********** SINGLE ELEVATOR FSM *****************************************
		case a := <-buttons:
			switch elevator.Behaviour {
			case elev.EB_DoorOpen:
				if requests.ShouldClearImmediately(elevator, a.Floor, a.Button) {
					openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
				} else {
					if a.Button == elevio.BT_Cab {
						elevator.Requests[a.Floor][a.Button] = true
						SendRequestsToBackup(elevator, proto, addr, backupFilePath)
					} else {
						worldView.Orders[a.Floor][a.Button].Order = elev.Unassigned
					}
				}
				wvMsg := wv.WorldViewMsg{
					Orders:        worldView.Orders,
					ID:            elevator.ID,
					ElevatorState: elevator}
				worldViewTx <- wvMsg

			case elev.EB_Moving:
				faultTimer.Reset(elev.FAULT_TIMEOUT)
				if a.Button == elevio.BT_Cab {
					elevator.Requests[a.Floor][a.Button] = true
					SendRequestsToBackup(elevator, proto, addr, backupFilePath)
				} else {
					worldView.Orders[a.Floor][a.Button].Order = elev.Unassigned
				}
				wvMsg := wv.WorldViewMsg{
					Orders:        worldView.Orders,
					ID:            elevator.ID,
					ElevatorState: elevator}
				worldViewTx <- wvMsg

			case elev.EB_Idle:
				if a.Button == elevio.BT_Cab {
					elevator.Requests[a.Floor][a.Button] = true
					SendRequestsToBackup(elevator, proto, addr, backupFilePath)
				} else {
					worldView.Orders[a.Floor][a.Button].Order = elev.Unassigned
				}

				pair := requests.ChooseDirection(elevator)
				elevator.Dir = pair.Dir
				elevator.Behaviour = pair.Behaviour

				switch elevator.Behaviour {
				case elev.EB_DoorOpen:
					elevio.SetDoorOpenLamp(true)
					openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
					elevator = requests.ClearAtCurrentFloor(elevator)
					SendRequestsToBackup(elevator, proto, addr, backupFilePath)
					faultTimer.Reset(elev.FAULT_TIMEOUT)
				case elev.EB_Moving:
					elevio.SetMotorDirection(elevator.Dir)
					faultTimer.Reset(elev.FAULT_TIMEOUT)

				case elev.EB_Idle:
					faultTimer.Reset(elev.FAULT_TIMEOUT)
				}

				wvMsg := wv.WorldViewMsg{
					Orders:        worldView.Orders,
					ID:            elevator.ID,
					ElevatorState: elevator}
				worldViewTx <- wvMsg
			}
			SetAllLights(worldView, elevator)

		case elevator.Floor = <-floors:
			elevio.SetFloorIndicator(elevator.Floor)

			switch elevator.Behaviour {
			case elev.EB_Moving:
				if requests.ShouldStop(elevator) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					elevator = requests.ClearAtCurrentFloor(elevator)
					worldView.Orders = wv.Orders_clearAtCurrentFloor(worldView, elevator).Orders
					openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
					elevator.Behaviour = elev.EB_DoorOpen
				}
				wvMsg := wv.WorldViewMsg{
					Orders:        worldView.Orders,
					ID:            elevator.ID,
					ElevatorState: elevator,
				}
				worldViewTx <- wvMsg
			}

		case a := <-obstruction:
			elevator.Obstructed = a
			wvMsg := wv.WorldViewMsg{
				Orders:        worldView.Orders,
				ID:            elevator.ID,
				ElevatorState: elevator,
				Fault:         a}
			worldViewTx <- wvMsg

		case a := <-stop:
			elevio.SetStopLamp(a)
			elevio.SetMotorDirection(elevio.MD_Stop)
			wvMsg := wv.WorldViewMsg{
				Orders:        worldView.Orders,
				ID:            elevator.ID,
				ElevatorState: elevator,
				Fault:         a}
			worldViewTx <- wvMsg

		case <-openDoorTimer.C:
			fmt.Println("Doortimer has timed out")

			if elevator.Obstructed {
				fmt.Println("Still obstructed")
				openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
				break
			}

			switch elevator.Behaviour {
			case elev.EB_DoorOpen:
				pair := requests.ChooseDirection(elevator)
				elevator.Dir = pair.Dir
				elevator.Behaviour = pair.Behaviour
				fmt.Print("New direction: ")
				fmt.Println(elevator.Dir)
				fmt.Print("New behaviour: ")
				fmt.Println(elevator.Behaviour)

				switch elevator.Behaviour {
				case elev.EB_DoorOpen:
					openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
					elevio.SetDoorOpenLamp(true)
					elevator = requests.ClearAtCurrentFloor(elevator)
					SendRequestsToBackup(elevator, proto, addr, backupFilePath)
					//SetAllLights(elevator)
				case elev.EB_Moving:
					elevio.SetDoorOpenLamp(false)
					elevio.SetMotorDirection(elevator.Dir)
				case elev.EB_Idle:
					elevio.SetDoorOpenLamp(false)
					elevio.SetMotorDirection(elevator.Dir)
				}
			}

		case <-faultTimer.C:

			if AreOrdersEmpty(worldView, elevator) {
				faultTimer.Reset(elev.FAULT_TIMEOUT)
			} else {
				pid := strconv.Itoa(os.Getpid())
				exec.Command("gnome-terminal", "--", "kill", "-TERM", pid).Run()
				wvMsg := wv.WorldViewMsg{Orders: worldView.Orders,
					ID:            elevator.ID,
					ElevatorState: elevator,
					Fault:         true}
				worldViewTx <- wvMsg
			}
		}
	}
}

func AreOrdersEmpty(worldView wv.WorldView, e elev.Elevator) bool {
	for i := 0; i < elev.N_FLOORS; i++ {
		for j := 0; j < elev.N_BUTTONS-1; j++ {
			if worldView.Orders[i][j].Order > elev.Completed {
				return false
			}
		}
	}
	for i := 0; i < elev.N_FLOORS; i++ {
		if e.Requests[i][elevio.BT_Cab] {
			return false
		}
	}
	return true
}

func SetAllLights(worldView wv.WorldView, e elev.Elevator) {
	for floor := 0; floor < elev.N_FLOORS; floor++ {
		for btn := 0; btn < elev.N_BUTTONS-1; btn++ {
			if worldView.Orders[floor][btn].Order == elev.Assigned {
				elevio.SetButtonLamp(elevio.ButtonType(btn), floor, true)
			} else {
				elevio.SetButtonLamp(elevio.ButtonType(btn), floor, false)
			}
		}
	}
	for floor := 0; floor < elev.N_FLOORS; floor++ {
		elevio.SetButtonLamp(elevio.ButtonType(elevio.BT_Cab), floor, e.Requests[floor][elevio.BT_Cab])

	}
}

func SendRequestsToBackup(e elev.Elevator, proto string, addr string, backupFilePath string) {

	var cabOrder = []byte{0, 0, 0, 0}

	for i := 0; i < elev.N_FLOORS; i++ {
		if e.Requests[i][2] {
			cabOrder[i] = 1
		} else {
			cabOrder[i] = 0
		}
	}

	os.WriteFile(backupFilePath, cabOrder, 0644)
}
