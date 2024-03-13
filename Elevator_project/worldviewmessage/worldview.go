package worldviewmessage

import elev "Elevator_project/elevator"

type WorldViewMsg struct {
	Orders        [elev.N_FLOORS][elev.N_BUTTONS - 1]int
	ID            string
	Fault         bool
	ElevatorState elev.Elevator
}

type WorldView struct {
	Orders     [elev.N_FLOORS][elev.N_BUTTONS - 1]int // 0 unassigned, 1 assigned, 2 completed
	Elevators  map[string]elev.Elevator
}

// func UpdateOrders(wv Request){

// 	for i := 0; i < elev.N_FLOORS; i++ {
// 		for j := 0; j < elev.N_BUTTONS; j++ {

// 		}
// 	}
// 	switch message.Queue {
// 	case Queue.Request == 0:

// 	}
// }

// func UpdateRequests(elevator){
// 	for i := 0; i < elev.N_FLOORS; i++ {
// 		for j := 0; j < elev.N_BUTTONS-1; j++ {
// 			// If mine add to local queue
// 			// elevator.Requests[a.Floor][a.Button] = true
// 		}
// 	}
// }
