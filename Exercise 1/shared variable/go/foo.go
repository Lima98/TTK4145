// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
	"time"
)

var i = 0

// func incrementing() {
//     //TODO: increment i 1000000 times
//     for j := 0; j < 1000000; j++ {
//         i++;
// 	}

// }

// func decrementing() {
//     //TODO: decrement i 1000000 times
//     for j := 0; j < 1000000; j++ {
//         i--;
//     }
// }
 func incrementing(ch chan int){
    i++
 }


func main() {
    // What does GOMAXPROCS do? What happens if you set it to 1?
    //GOMAXPROCS setter maks antall prosesser (her tråder) som kan kjøres samtidig. Sette til 1--> kan bare kjøre en tråd 
    runtime.GOMAXPROCS(2)
    channel_1 := make(chan int)
    channel_2 := make(chan int)


    // go incrementing()
    // go decrementing()


    go incrementing(channel_1)
    go decrementing(channel_2)


    select{
    case new_i := <- channel_1:
        incrementing()
        print("Incremented i")
    case ch <- i--:
        print("Decremented i")
    default:
        print "no communication"
    }

	
    // TODO: Spawn both functions as goroutines
    //Goroutine er basically en thread i Go, altså det som kaller på funksjonen, og får funksjonen til å kjøre

   

	
    // We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
    // We will do it properly with channels soon. For now: Sleep.
    time.Sleep(500*time.Millisecond)
    Println("The magic number is:", i)
}
