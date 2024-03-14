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
			worldView.Orders[i][j].Order = elev.Completed
			worldView.Orders[i][j].ElevatorsThatKnow = make(map[string]bool)
		}
	} // legg inn sjekk for a kun riktig heis kan sette til completed

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

	go SendAllTheTime(worldView, elevator, worldViewTx)

	for {
		fmt.Println("---------------------------------------------------------------")
		fmt.Println("---------------------------------------------------------------")
		requests.PrintRequests(elevator)
		elev.PrintBehaviour(elevator)

		select {
		// NETWORK
		case a := <-peerUpdateCh:
			fmt.Println("PEER UPDATE")
			fmt.Println(a) // DEtte må vi få sett på mtp. orderefordeling
			fmt.Println("-")
			peerList = a.Peers

		case a := <-worldViewRx:
			worldView.Elevators[a.ID] = a.ElevatorState //Burde ha noe kontroll?

			fmt.Print("WORLDVIEW RECEIVED: ")
			fmt.Println(a)
			//if a.ID == id {break}
			for i := 0; i < elev.N_FLOORS; i++ {
				for j := 0; j < elev.N_BUTTONS-1; j++ {
					switch worldView.Orders[i][j].Order {
					case elev.Completed:
						switch a.Orders[i][j].Order {
						case elev.Unassigned:
							worldView.Orders[i][j].Order = a.Orders[i][j].Order
							//worldView.Orders[i][j].ElevatorsThatKnow = a.Orders[i][j].ElevatorsThatKnow
							fmt.Println("Elevator " + id + " updated worldview from Completed to Unassigned.")
						case elev.Assigned:
							worldView.Orders[i][j].Order = a.Orders[i][j].Order
							worldView.Orders[i][j].ElevatorsThatKnow = a.Orders[i][j].ElevatorsThatKnow
							fmt.Println("Elevator " + id + " updated worldview from Completed to Assigned.")
						case elev.Completed:
						}
					case elev.Unassigned:
						switch a.Orders[i][j].Order {
						case elev.Unassigned:
						case elev.Assigned:
							worldView.Orders[i][j].Order = a.Orders[i][j].Order
							worldView.Orders[i][j].ElevatorsThatKnow = a.Orders[i][j].ElevatorsThatKnow
							fmt.Println("Elevator " + id + " updated worldview from Unassigned to Assigned.")
						case elev.Completed:
						}
					case elev.Assigned:
						switch a.Orders[i][j].Order {
						case elev.Unassigned:
						case elev.Assigned:
						case elev.Completed:
							for k := 0; k < len(peerList); k++ {
								if !worldView.Orders[i][j].ElevatorsThatKnow[peerList[k]] {
									fmt.Println("entered break if statement")
									break
								}
							}
							worldView.Orders[i][j].Order = a.Orders[i][j].Order
							fmt.Println("Elevator " + id + " updated worldview from Assigned to Completed.")

							worldView.Orders[i][j].ElevatorsThatKnow = make(map[string]bool)
							fmt.Println("Created new map empty  ")
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

			// Må endres når hall assigner er implementert
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

		case a := <-buttons:
			switch elevator.Behaviour {
			case elev.EB_DoorOpen:
				if requests.Requests_shouldClearImmediately(elevator, a.Floor, a.Button) {
					openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
				} else {
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
				}

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

				pair := requests.Requests_chooseDirection(elevator)
				elevator.Dir = pair.Dir
				elevator.Behaviour = pair.Behaviour

				switch elevator.Behaviour {
				case elev.EB_DoorOpen:
					elevio.SetDoorOpenLamp(true)
					openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
					elevator = requests.Requests_clearAtCurrentFloor(elevator)
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

		case elevator.Floor = <-floors:
			elevio.SetFloorIndicator(elevator.Floor)

			switch elevator.Behaviour {
			case elev.EB_Moving:
				if requests.Requests_shouldStop(elevator) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					elevator = requests.Requests_clearAtCurrentFloor(elevator)
					worldView.Orders = wv.Orders_clearAtCurrentFloor(worldView, elevator).Orders
					openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
					//SetAllLights(elevator)
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
				exec.Command("gnome-terminal", "--", "kill", "-TERM", pid).Run() //opens a new window so might be messy
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
}

func AreOrdersEmpty(worldView wv.WorldView, e elev.Elevator) bool {
	for i := 0; i < elev.N_FLOORS; i++ {
		for j := 0; j < elev.N_BUTTONS-1; j++ {
			if worldView.Orders[i][j].Order < elev.Completed {
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

func SendAllTheTime(worldView wv.WorldView, elevator elev.Elevator, worldViewTx chan wv.WorldViewMsg) {
	time.Sleep(100 * time.Millisecond)
	wvMsg := wv.WorldViewMsg{
		Orders:        worldView.Orders,
		ID:            elevator.ID,
		ElevatorState: elevator,
	}
	worldViewTx <- wvMsg
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

	var cabOrder = []byte{0, 0, 0, 0} // Kan vi gjøre dette basert på numfloors?

	for i := 0; i < elev.N_FLOORS; i++ {
		if e.Requests[i][2] {
			cabOrder[i] = 1
		} else {
			cabOrder[i] = 0
		}
	}

	os.WriteFile(backupFilePath, cabOrder, 0644)
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
