// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
	"time"
)

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
 func incrementing(ch chan int, quit chan int){
    for j := 0; j < 1000000; j++ {
        ch <- 1        
  }
    quit<-1
 }

 func decrementing(ch chan int, quit chan int){
    for j := 0; j < 1000000; j++ {
        ch <- 1        
  }
    quit<-1
 }


 func server(ch_increment chan int, ch_decrement chan int, ch_read chan int){
    var i = 0
    for {
    select{
    case <-ch_increment:
        i++
    case <- ch_decrement:
        i--
    case ch_read <- i:

    }
}
}

func main() {
    // What does GOMAXPROCS do? What happens if you set it to 1?
    //GOMAXPROCS setter maks antall prosesser (her tråder) som kan kjøres samtidig. Sette til 1--> kan bare kjøre en tråd 
    runtime.GOMAXPROCS(2)
    ch_increment := make(chan int)
    ch_decrement:= make(chan int)
    ch_read:= make(chan int)
    ch_quit:= make(chan int)



    // go incrementing()
    // go decrementing()

    go server(ch_increment,ch_decrement,ch_read)
    go incrementing(ch_increment, ch_quit)
    go decrementing(ch_decrement, ch_quit)


    <-ch_quit
    <-ch_quit
	
    // TODO: Spawn both functions as goroutines
    //Goroutine er basically en thread i Go, altså det som kaller på funksjonen, og får funksjonen til å kjøre

   

	
    // We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
    // We will do it properly with channels soon. For now: Sleep.
    time.Sleep(500*time.Millisecond)
    Println("The magic number is:", <- ch_read)
}
