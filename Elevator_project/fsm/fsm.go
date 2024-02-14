package fsm

import (
	elevio "Elevator_project/driver-go/elevio"
	elev "Elevator_project/elevator"
	"fmt"
)

var elevator = elev.Elevator{Floor: -1,
                        Dir: elevio.MD_Stop,
                        Behaviour: elev.EB_Idle}


func Fsm_test(btn_floor int, btn_type elevio.ButtonType){
    fmt.Println(btn_floor)
	fmt.Println(elevator.Floor)
}

func fsm_onRequestButtonPress(btn_floor int, btn_type elevio.ButtonType){
    // printf("\n\n%s(%d, %s)\n", __FUNCTION__, btn_floor, elevio_button_toString(btn_type));
    // elevator_print(elevator);
    
    switch(elevator.Behaviour){
    case elev.EB_DoorOpen:
        if(requests_shouldClearImmediately(elevator, btn_floor, btn_type)){
            timer_start(elevator.config.doorOpenDuration_s);
        } else {
            elevator.requests[btn_floor][btn_type] = 1;
        }
        break;

    case EB_Moving:
        elevator.requests[btn_floor][btn_type] = 1;
        break;
        
    case EB_Idle:    
        elevator.requests[btn_floor][btn_type] = 1;
        pair := req.Requests_chooseDirection(elevator);
        elevator.dirn = pair.dirn;
        elevator.behaviour = pair.behaviour;
        switch(pair.behaviour){
        case EB_DoorOpen:
            outputDevice.doorLight(1);
            timer_start(elevator.config.doorOpenDuration_s);
            elevator = requests_clearAtCurrentFloor(elevator);
            break;

        case EB_Moving:
            outputDevice.motorDirection(elevator.dirn);
            break;
            
        case EB_Idle:
            break;
        }
        break;
    }
    
    setAllLights(elevator);
    
    printf("\nNew state:\n");
    elevator_print(elevator);
}
