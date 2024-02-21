// package fsm

// import (
// 	elevio "Elevator_project/driver-go/elevio"
// 	elev "Elevator_project/elevator"
// 	"Elevator_project/requests"
// 	"fmt"
// 	"time"
// )



// func Statemachine(){

//     buttons := make(chan elevio.ButtonEvent)
//     floor  := make(chan int)
//     obstruction   := make(chan bool)
//     stop    := make(chan bool)    
    
//     go elevio.PollButtons(buttons)
//     go elevio.PollFloorSensor(floor)
//     go elevio.PollObstructionSwitch(obstruction)
//     go elevio.PollStopButton(stop)

//     var elevator = elev.Elevator{Floor: 1,
//                                  Dir: elevio.MD_Stop,
//                                  Behaviour: elev.EB_Idle,
//                                  Obstructed: false}
                                 
//     openDoorTimer := time.NewTimer(1000*time.Second)
    
//     // Initialisering, den går alltid til 1. etasje uansett. 
//     // Vi vil kun flytte oss hvis vi er mellom etasjer.
//     // Bør fikses
//     elevio.SetMotorDirection(elevio.MD_Down);
//     elevator.Dir = elevio.MD_Down;
//     elevator.Behaviour = elev.EB_Moving;
//     elevator.Requests[0][2] = true

//     for {
//         requests.PrintRequests(elevator)
//         elev.PrintBehaviour(elevator)

//         switch(elevator.Behaviour){
//             case elev.EB_DoorOpen:
//                 select{
//                     case a := <- buttons:
//                         elevator.Requests[a.Floor][a.Button] = true;
//                     case a := <- floor:

//                     case elevator.Obstructed  = <- obstruction:  

//                     case a := <- stop:

//                     case <- openDoorTimer.C:

//                     }
//             case elev.EB_Idle:
//                 select{
//                     case a := <- buttons:
//                         elevator.Requests[a.Floor][a.Button] = true;
//                     case a := <- floor:
//                     case elevator.Obstructed  = <- obstruction:  
//                     case a := <- stop:
//                     case <- openDoorTimer.C:
//                 }
//             case elev.EB_Moving:
//                 select{
//                     case a := <- buttons:
//                         elevator.Requests[a.Floor][a.Button] = true;
//                     case a := <- floor:
//                     case elevator.Obstructed  = <- obstruction:  
//                     case a := <- stop:
//                     case <- openDoorTimer.C:
//                 }

//         }
//     }
// }


// func SetAllLights(e elev.Elevator){
// 	for floor := 0; floor < elev.N_FLOORS; floor++ {
//         for btn := 0; btn < elev.N_BUTTONS; btn++{
//             elevio.SetButtonLamp(elevio.ButtonType(btn), floor, e.Requests[floor][btn]);
//         }
//     }
// }

