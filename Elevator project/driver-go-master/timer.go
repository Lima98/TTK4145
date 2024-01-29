package main

import (
	"fmt"
	"time"
)

type Timer struct {
	endTime  time.Time
	active   bool
	duration time.Duration
}

func newTimer() *Timer {
	return &Timer{}
}

func (t *Timer) start(duration time.Duration) {
	t.endTime = time.Now().Add(duration)
	t.active = true
	t.duration = duration
}

func (t *Timer) stop() {
	t.active = false
}

func (t *Timer) timedOut() bool {
	return t.active && time.Now().After(t.endTime)
}

func main() {
	fmt.Println("Started!")

	timer := newTimer()

	// Example usage:
	timer.start(5 * time.Second)

	for {
		if timer.timedOut() {
			fmt.Println("Timer timed out!")
			timer.stop()
		}

		// Do other tasks or sleep
		time.Sleep(100 * time.Millisecond)
	}
}
