package main

import (
	"Elevator_project/driver-go/elevio"
	"Elevator_project/fsm"
	"Elevator_project/timer"
)



func main(){

    numFloors := 4

    elevio.Init("localhost:15657", numFloors)
    
    var d elevio.MotorDirection = elevio.MD_Stop
    elevio.SetMotorDirection(d)
    
    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)    
    
    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)
    //masse masse kode

    // if  <- drv_floors == -1 {
    //     fsm.Fsm_onInitBetweenFloors()
    // } // FIKS DETTE SENERE NÅR DERE SKJØNNE HVA FAEN SOM FOREGÅR


    for {
        select{
        case a := <- drv_buttons:
            fsm.Fsm_onRequestButtonPress(a.Floor,a.Button)

        case a := <- drv_floors:
            elevio.SetFloorIndicator(a)
            fsm.Fsm_onFloorArrival(a)

        case a := <- drv_obstr:
            elevio.SetStopLamp(a)
            elevio.SetMotorDirection(elevio.MD_Stop)
        
        case a := <- drv_stop:
            elevio.SetStopLamp(a)
            elevio.SetMotorDirection(elevio.MD_Stop)

        }
        if(timer.Timer_timeout()){
            timer.Timer_stop()
            fsm.Fsm_onDoorTimeout()
        }
    }    
}
