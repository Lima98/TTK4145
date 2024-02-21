package fsm

import (
	elevio "Elevator_project/driver-go/elevio"
	elev "Elevator_project/elevator"
	"Elevator_project/requests"
	"fmt"
	"time"
)



func Statemachine(){

    buttons := make(chan elevio.ButtonEvent)
    floor  := make(chan int)
    obstruction   := make(chan bool)
    stop    := make(chan bool)    
    
    go elevio.PollButtons(buttons)
    go elevio.PollFloorSensor(floor)
    go elevio.PollObstructionSwitch(obstruction)
    go elevio.PollStopButton(stop)

    var elevator = elev.Elevator{Floor: 1,
                                 Dir: elevio.MD_Stop,
                                 Behaviour: elev.EB_Idle,
                                 Obstructed: false}
                                 
    openDoorTimer := time.NewTimer(1000*time.Second)
    
    // Initialisering, den går alltid til 1. etasje uansett. 
    // Vi vil kun flytte oss hvis vi er mellom etasjer.
    // Bør fikses
    elevio.SetMotorDirection(elevio.MD_Down);
    elevator.Dir = elevio.MD_Down;
    elevator.Behaviour = elev.EB_Moving;
    elevator.Requests[0][2] = true

    for {
        requests.PrintRequests(elevator)
        select{
        case a := <- buttons:
            switch(elevator.Behaviour){
                case elev.EB_DoorOpen:
                    if(requests.Requests_shouldClearImmediately(elevator, a.Floor, a.Button)){
                        openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
                    } else {
                        elevator.Requests[a.Floor][a.Button] = true;
                    }
            
                case elev.EB_Moving:
                    elevator.Requests[a.Floor][a.Button] = true;
                    
                case elev.EB_Idle:    
                    elevator.Requests[a.Floor][a.Button] = true;
                    pair := requests.Requests_chooseDirection(elevator);
                    elevator.Dir = pair.Dir;
                    elevator.Behaviour = pair.Behaviour;
                    
                    switch(pair.Behaviour){
                        case elev.EB_DoorOpen:
                            elevio.SetDoorOpenLamp(true)
                            openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
                            elevator = requests.Requests_clearAtCurrentFloor(elevator);
                
                        case elev.EB_Moving:
                            elevio.SetMotorDirection(elevator.Dir)
                            
                        case elev.EB_Idle:

                    }
            }
            
            SetAllLights(elevator); // skal kanskje bort når vi har et bestillingssystem

        case a := <- floor:
            elevio.SetFloorIndicator(a)
            fmt.Println("OnFloorArrival")

            elevator.Floor = a;

            
            switch(elevator.Behaviour){
            case elev.EB_Moving:
                if(requests.Requests_shouldStop(elevator)){
                    fmt.Println("shouldstop")
                    elevio.SetMotorDirection(elevio.MD_Stop)
                    elevio.SetDoorOpenLamp(true)
                    elevator = requests.Requests_clearAtCurrentFloor(elevator)
                    openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
                    SetAllLights(elevator)
                    elevator.Behaviour = elev.EB_DoorOpen;
                }
        }

        case a := <- obstruction:  
            elevator.Obstructed = a   
        
        case a := <- stop:
            elevio.SetStopLamp(a)
            elevio.SetMotorDirection(elevio.MD_Stop)

        case <- openDoorTimer.C:
            fmt.Println("timer has timed out")
            fmt.Println("OnDoorTimeout")
            fmt.Println(elevator)
        
            if elevator.Obstructed {
                fmt.Println("Still obstructed")
                openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
                return
            }
        
            elevio.SetDoorOpenLamp(false)
        
        
            switch(elevator.Behaviour){
            case elev.EB_DoorOpen:
                pair := requests.Requests_chooseDirection(elevator)
                elevator.Dir = pair.Dir;
                elevator.Behaviour = pair.Behaviour;
        
        
                switch(elevator.Behaviour){
                    case elev.EB_DoorOpen:
                        openDoorTimer.Reset(elev.OPEN_DOOR_TIME)
                        elevator = requests.Requests_clearAtCurrentFloor(elevator);
                        SetAllLights(elevator)
                    case elev.EB_Moving:
                    case elev.EB_Idle:
                        elevio.SetDoorOpenLamp(true)
                        elevio.SetMotorDirection(elevator.Dir)
        
                }
        
            }
        }
    }  
}


func SetAllLights(e elev.Elevator){
	for floor := 0; floor < elev.N_FLOORS; floor++ {
        for btn := 0; btn < elev.N_BUTTONS; btn++{
            elevio.SetButtonLamp(elevio.ButtonType(btn), floor, e.Requests[floor][btn]);
        }
    }
}

