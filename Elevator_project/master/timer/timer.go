package timer

import (
	"fmt"
	"time"
)

// viktig spørsmål!! går det an å bare bruke en sleep som timer? lol

// oversatt c-koden. ikke satt meg inn i disse timer-funksjonene i go, de kan det jon hende man heller vil bruke.
// evt slo de meg at dette er farlig likt å bare sleepe.. ref spm over

var Timer_active = false
var Timer_end_time time.Time

func Timer_start(num_seconds float32) {
	current_time := time.Now()
	Timer_end_time = current_time.Add(time.Duration(num_seconds * float32(time.Second)))
	Timer_active = true
}

func Timer_stop() {
	Timer_active = false
}

func Timer_timeout() bool {
	current_time := time.Now()
	return (current_time.After(Timer_end_time))
}

// run_timer er mest en testfunksjon for å sjekke at start og stop funker som de skal:)
func Run_timer(num_seconds float32) {
	Timer_start(num_seconds)
	fmt.Println("started")
	for {
		if Timer_timeout() {
			Timer_stop()
			fmt.Println("stopped")
			break
		}
	}
}
