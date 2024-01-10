
// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
	"time"
)

var i = 0
var ch1, ch2, ch3 chan int
var inc, dec, done int


func incrementing() {
    //TODO: increment i 1000000 times
    for j := 0; j < 1000000; j++ {
        ch1 <= inc
	}
    ch2 <= done
}

func decrementing() {
    //TODO: decrement i 1000000 times
    for j := 0; j < 1000000; j++ {
        ch1 <= dec
    }
    ch3 <= done
}




func main() {
    // What does GOMAXPROCS do? What happens if you set it to 1?
    //GOMAXPROCS setter maks antall prosesser (her tråder) som kan kjøres samtidig. Sette til 1--> kan bare kjøre en tråd 
    runtime.GOMAXPROCS(2)    

    // TODO: Spawn both functions as goroutines
    //Goroutine er basically en thread i Go, altså det som kaller på funksjonen, og får funksjonen til å kjøre

    go incrementing();
    go decrementing();

    for (done <= ch2 && done <= ch3 == 0) {
        select{
        case inc <= ch1:
            i++
        case dec <= ch1:
            i--
        }
    }
	
    // We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
    // We will do it properly with channels soon. For now: Sleep.
    time.Sleep(500*time.Millisecond)
    Println("The magic number is:", i)
}
