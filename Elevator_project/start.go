package main

import (
	"fmt"
	"os/exec"
	"strconv"
)

func main() {
	port := 15657

	number_of_elevators := 3

	for i := 0; i < number_of_elevators; i++ {
		cmd_simulator := exec.Command("gnome-terminal", "--", "./SimElevatorServer", "--port", strconv.Itoa(port+i))
		err := cmd_simulator.Run()
		if err != nil {
			fmt.Println("Error starting: ", err)
		}

		cmd_main := exec.Command("gnome-terminal", "--", "go", "run", "main.go", "-id="+strconv.Itoa(i))
		err = cmd_main.Run()
		if err != nil {
			fmt.Println("Error starting: ", err)
		}
	}

}
