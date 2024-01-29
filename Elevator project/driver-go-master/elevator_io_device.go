package main

// ElevInputDevice represents the input device of the elevator.
type ElevInputDevice struct {
	FloorSensor   func() int
	RequestButton func(int, Button) int
	StopButton    func() int
	Obstruction   func() int
}

// ElevOutputDevice represents the output device of the elevator.
type ElevOutputDevice struct {
	FloorIndicator     func(int)
	RequestButtonLight func(int, Button, int)
	DoorLight          func(int)
	StopButtonLight    func(int)
	MotorDirection     func(Dirn)
}

// DirnToString converts a Dirn to a string.
func DirnToString(d Dirn) string {
	switch d {
	case D_Up:
		return "D_Up"
	case D_Down:
		return "D_Down"
	case D_Stop:
		return "D_Stop"
	default:
		return "D_UNDEFINED"
	}
}

// ButtonToString converts a Button to a string.
func ButtonToString(b Button) string {
	switch b {
	case B_HallUp:
		return "B_HallUp"
	case B_HallDown:
		return "B_HallDown"
	case B_Cab:
		return "B_Cab"
	default:
		return "B_UNDEFINED"
	}
}

// GetInputDevice returns the input device of the elevator.
func GetInputDevice() ElevInputDevice {
	return ElevInputDevice{
		FloorSensor:   elevatorHardwareGetFloorSensorSignal,
		RequestButton: wrapRequestButton,
		StopButton:    elevatorHardwareGetStopSignal,
		Obstruction:   elevatorHardwareGetObstructionSignal,
	}
}

// GetOutputDevice returns the output device of the elevator.
func GetOutputDevice() ElevOutputDevice {
	return ElevOutputDevice{
		FloorIndicator:     elevatorHardwareSetFloorIndicator,
		RequestButtonLight: wrapRequestButtonLight,
		DoorLight:          elevatorHardwareSetDoorOpenLamp,
		StopButtonLight:    elevatorHardwareSetStopLamp,
		MotorDirection:     wrapMotorDirection,
	}
}

func wrapRequestButton(f int, b Button) int {
	return elevatorHardwareGetButtonSignal(int(b), f)
}

func wrapRequestButtonLight(f int, b Button, v int) {
	elevatorHardwareSetButtonLamp(int(b), f, v)
}

func wrapMotorDirection(d Dirn) {
	elevatorHardwareSetMotorDirection(int(d))
}
