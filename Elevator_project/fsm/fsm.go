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
	"net"
	"os"
	"strconv"
	"time"
)

func Statemachine(proto string, addr string, cabOrders []byte, id string) {

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

	idToString := make(map[string]string)
	idToString["0"] = "one"
	idToString["1"] = "two"
	idToString["2"] = "three"

	var elevator = elev.Elevator{Floor: 1,
		Dir:        elevio.MD_Stop,
		Behaviour:  elev.EB_Idle,
		Obstructed: false,
		ID:         id}

	var worldView = wv.WorldView{}
	worldView.Elevators = make(map[string]elev.Elevator)
	var peerList []string

	for i := 0; i < elev.N_FLOORS; i++ {
		for j := 0; j < elev.N_BUTTONS-1; j++ {
			worldView.Orders[i][j] = 2
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
	faultTimer := time.NewTimer(elev.FAULT_TIMEOUT)

	select {
	case <-floors:
	default:
		elevio.SetMotorDirection(elevio.MD_Down)
		elevator.Dir = elevio.MD_Down
		elevator.Behaviour = elev.EB_Moving
	}

	for {
		fmt.Println("---------------------------------------------------------------")
		fmt.Println("---------------------------------------------------------------")
		requests.PrintRequests(elevator)
		elev.PrintBehaviour(elevator)

		for i := 0; i < elev.N_FLOORS; i++ {
			for j := 0; j < elev.N_BUTTONS-1; j++ {
				if !elevator.Requests[i][j] {
					worldView.Orders[i][j] = 2
				}
			}
		}

		SendRequestsToBackup(elevator, proto, addr)
		// Need to send the queue to the master queue
		// Need to send and recieve the queue on the network

		select {
		// NETWORK
		case a := <-peerUpdateCh:
			fmt.Println("PEER UPDATE")
			fmt.Println(a) // DEtte må vi få sett på mtp. orderefordeling
			fmt.Println("-")
			peerList = a.Peers

		case a := <-worldViewRx:
			fmt.Print("WORLDVIEW RECEIVED: ")
			fmt.Println(a)
			if a.ID == id {break}
			for i := 0; i < elev.N_FLOORS; i++ {
				for j := 0; j < elev.N_BUTTONS-1; j++ {
					switch worldView.Orders[i][j] {
					case 0:
						switch a.Orders[i][j] {
							case 0:
							case 1: worldView.Orders[i][j] = a.Orders[i][j]
							case 2:
						}
					case 1:
						switch a.Orders[i][j] {
							case 0:
							case 1:
							case 2:
								for i := 0; i < elev.N_ELEVATORS; i++ {
									if worldView.Elevators[idToString[strconv.Itoa(i)]].Requests[i][j] {
										worldView.Orders[i][j] = a.Orders[i][j]
									}
								} 
					}
					case 2:
						switch a.Orders[i][j] {
							case 0: worldView.Orders[i][j] = a.Orders[i][j]
							case 1:
							case 2:
					}
					}
				}
			}
			hallAssignments := hra.HallRequestAssigner(worldView.Orders, worldView.Elevators, peerList)

			fmt.Println("HALL ASSIGNMENTS:")
			fmt.Println(hallAssignments["one"])

			// Må endres når hall assigner er implementert
			for i := 0; i < elev.N_FLOORS; i++ {
				for j := 0; j < elev.N_BUTTONS-1; j++ {
					elevator.Requests[i][j] = hallAssignments[idToString[elevator.ID]][i][j]
					if hallAssignments[idToString[elevator.ID]][i][j] {
						worldView.Orders[i][j] = 1
					}
				}
			}
		

		case a := <-buttons:
			switch elevator.Behaviour {
			case elev.EB_DoorOpen:
				if requests.Requests_shouldClearImmediately(elevator, a.Floor, a.Button) {
					openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
				} else {
					if a.Button == elevio.BT_Cab {
						elevator.Requests[a.Floor][a.Button] = true
					} else {
						worldView.Orders[a.Floor][a.Button] = 0
					}
					wvMsg := wv.WorldViewMsg{Orders: worldView.Orders,
						ID:            elevator.ID,
						ElevatorState: elevator}
					worldViewTx <- wvMsg
				}

			case elev.EB_Moving:
				faultTimer.Reset(elev.FAULT_TIMEOUT)
				if a.Button == elevio.BT_Cab {
					elevator.Requests[a.Floor][a.Button] = true
				} else {
					worldView.Orders[a.Floor][a.Button] = 0
				}
				wvMsg := wv.WorldViewMsg{Orders: worldView.Orders,
					ID:            elevator.ID,
					ElevatorState: elevator}
				worldViewTx <- wvMsg

			case elev.EB_Idle:
				if a.Button == elevio.BT_Cab {
					elevator.Requests[a.Floor][a.Button] = true
				} else {
					worldView.Orders[a.Floor][a.Button] = 0
				}
				wvMsg := wv.WorldViewMsg{Orders: worldView.Orders,
					ID:            elevator.ID,
					ElevatorState: elevator}
				worldViewTx <- wvMsg

				pair := requests.Requests_chooseDirection(elevator)
				elevator.Dir = pair.Dir
				elevator.Behaviour = pair.Behaviour

				switch elevator.Behaviour {
				case elev.EB_DoorOpen:
					elevio.SetDoorOpenLamp(true)
					openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
					elevator = requests.Requests_clearAtCurrentFloor(elevator)

				case elev.EB_Moving:
					elevio.SetMotorDirection(elevator.Dir)

				case elev.EB_Idle:
					//faultTimer.Reset(elev.FAULT_TIMEOUT)
				}
			}

			SetAllLights(elevator) // skal kanskje bort når vi har et bestillingssystem

		case elevator.Floor = <-floors:
			elevio.SetFloorIndicator(elevator.Floor)

			switch elevator.Behaviour {
			case elev.EB_Moving:
				if requests.Requests_shouldStop(elevator) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					elevator = requests.Requests_clearAtCurrentFloor(elevator)
					openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
					SetAllLights(elevator)
					elevator.Behaviour = elev.EB_DoorOpen
				}
			}

		case a := <-obstruction:
			wvMsg := wv.WorldViewMsg{Orders: worldView.Orders,
				ID:            elevator.ID,
				ElevatorState: elevator,
				Fault:         a}
			worldViewTx <- wvMsg
			elevator.Obstructed = a

		case a := <-stop:
			elevio.SetStopLamp(a)
			elevio.SetMotorDirection(elevio.MD_Stop)

		case <-openDoorTimer.C:
			fmt.Println("Doortimer has timed out")

			if elevator.Obstructed {
				fmt.Println("Still obstructed")
				openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
				break
			}

			switch elevator.Behaviour {
			case elev.EB_DoorOpen:
				pair := requests.Requests_chooseDirection(elevator)
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
					elevator = requests.Requests_clearAtCurrentFloor(elevator)
					SetAllLights(elevator)
				case elev.EB_Moving:
					elevio.SetDoorOpenLamp(false)
					elevio.SetMotorDirection(elevator.Dir)
				case elev.EB_Idle:
					elevio.SetDoorOpenLamp(false)
					elevio.SetMotorDirection(elevator.Dir)
				}
			}
		case <-faultTimer.C:
			fmt.Println("##\n##\n##\n##\n##\n##\n##\n##\n##\n##\n##\n##\n##")
			wvMsg := wv.WorldViewMsg{Orders: worldView.Orders,
				ID:            elevator.ID,
				ElevatorState: elevator,
				Fault:         true}
			worldViewTx <- wvMsg
			fmt.Println("Fault is set to", wvMsg.Fault)
			fmt.Println("##\n##\n##\n##\n##\n##\n##\n##\n##\n##\n##\n##\n##")

		}
	}
}

func SetAllLights(e elev.Elevator) {
	for floor := 0; floor < elev.N_FLOORS; floor++ {
		for btn := 0; btn < elev.N_BUTTONS; btn++ {
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, e.Requests[floor][btn])
		}
	}
}

func SendRequestsToBackup(e elev.Elevator, proto string, addr string) {

	var cabOrder = []byte{0, 0, 0, 0} // Kan vi gjøre dette basert på numfloors?

	for i := 0; i < elev.N_FLOORS; i++ {
		if e.Requests[i][2] {
			cabOrder[i] = 1
		} else {
			cabOrder[i] = 0
		}
	}

	os.WriteFile("./autorestart/cab_orders.txt", cabOrder, 0644)
}

func ReceiveRequestsFromBackup(e *elev.Elevator, proto string, addr string) {
	conn, err := net.ListenPacket(proto, addr)

	//conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 1024)
	num_of_bytes, _, _ := conn.ReadFrom(buf)
	fmt.Print("FSM receive from backup: ")
	fmt.Println(num_of_bytes)

	for i := 0; i < elev.N_FLOORS; i++ {
		if buf[i] == 1 {
			e.Requests[i][2] = true
		} else {
			e.Requests[i][2] = false
		}
	}
	if err == nil {
		conn.Close()
	}
}
