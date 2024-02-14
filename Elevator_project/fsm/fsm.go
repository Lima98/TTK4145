package fsm

import (
	elevio "Elevator_project/driver-go/elevio"
	elev "Elevator_project/elevator"
	"Elevator_project/requests"
	"Elevator_project/timer"
)

var elevator = elev.Elevator{Floor: 1,
                        Dir: elevio.MD_Stop,
                        Behaviour: elev.EB_Idle}


func Fsm_test(){
    requests.PrintRequests(elevator)
}

func Fsm_setAllLights(e elev.Elevator){
	for floor := 0; floor < elev.N_FLOORS; floor++ {
        for btn := 0; btn < elev.N_BUTTONS; btn++{
            elevio.SetButtonLamp(elevio.ButtonType(btn), floor, e.Requests[floor][btn]);
        }
    }
}

func Fsm_onInitBetweenFloors(){
    elevio.SetMotorDirection(elevio.MD_Down);
    elevator.Dir = elevio.MD_Down;
    elevator.Behaviour = elev.EB_Moving;
}

func Fsm_onRequestButtonPress(btn_floor int, btn_type elevio.ButtonType){
    // printf("\n\n%s(%d, %s)\n", __FUNCTION__, btn_floor, elevio_button_toString(btn_type));
    // elevator_print(elevator);
    
    switch(elevator.Behaviour){
    case elev.EB_DoorOpen:
        if(requests.Requests_shouldClearImmediately(elevator, btn_floor, btn_type)){
            timer.Timer_start(3); // maybe add a config or variable for door time
        } else {
            elevator.Requests[btn_floor][btn_type] = true;
        }
        break;

    case elev.EB_Moving:
        elevator.Requests[btn_floor][btn_type] = true;
        break;
        
    case elev.EB_Idle:    
        elevator.Requests[btn_floor][btn_type] = true;
        pair := requests.Requests_chooseDirection(elevator);
        elevator.Dir = pair.Dir;
        elevator.Behaviour = pair.Behaviour;
        switch(pair.Behaviour){
        case elev.EB_DoorOpen:
            elevio.SetDoorOpenLamp(true)
            timer.Timer_start(3);  // maybe add a config or variable for door time
            elevator = requests.Requests_clearAtCurrentFloor(elevator);
            break;

        case elev.EB_Moving:
            elevio.SetMotorDirection(elevator.Dir)
            break;
            
        case elev.EB_Idle:
            break;
        }
        break;
    }
    
    Fsm_setAllLights(elevator); 
}

 func Fsm_onFloorArrival(newFloor int){
    // printf("\n\n%s(%d)\n", __FUNCTION__, newFloor);
    // elevator_print(elevator);
    
    elevator.Floor = newFloor;

    
    switch(elevator.Behaviour){
    case elev.EB_Moving:
        if(requests.Requests_shouldStop(elevator)){
            elevio.SetMotorDirection(elevio.MD_Stop)
            elevio.SetDoorOpenLamp(true)
            elevator = requests.Requests_clearAtCurrentFloor(elevator)
            timer.Timer_start(3) // HUSK Å FIKSE DENENN ØDØRITMEMIER
            Fsm_setAllLights(elevator)
            elevator.Behaviour = elev.EB_DoorOpen;
        }
        break;
    default:
        break;
    }
 }

 func Fsm_onDoorTimeout(){
    // printf("\n\n%s()\n", __FUNCTION__);
    // elevator_print(elevator);
    
    switch(elevator.Behaviour){
    case elev.EB_DoorOpen:
        pair := requests.Requests_chooseDirection(elevator)
        elevator.Dir = pair.Dir;
        elevator.Behaviour = pair.Behaviour;
        
        switch(elevator.Behaviour){
        case elev.EB_DoorOpen:
            timer.Timer_start(3) // HUSK Å FIKS
            elevator = requests.Requests_clearAtCurrentFloor(elevator);
            Fsm_setAllLights(elevator)
            break;
        case elev.EB_Moving:
        case elev.EB_Idle:
            elevio.SetDoorOpenLamp(true)
            elevio.SetMotorDirection(elevator.Dir)
            break;
        }
        
        break;
    default:
        break;
    }
    
    // printf("\nNew state:\n");
    // elevator_print(elevator);
}
