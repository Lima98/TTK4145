package network

import (
	"Elevator_project/network/network/bcast"
	"Elevator_project/network/network/localip"
	"Elevator_project/network/network/peers"
	elev "Elevator_project/elevator"
	"flag"
	"fmt"
	"os"
	//"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.
type WorldViewMsg struct {
	Orders [elev.N_FLOORS][elev.N_BUTTONS-1]int
	ID string
}

func Network(worldViewTx chan WorldViewMsg, worldViewRx chan WorldViewMsg) {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)

	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// We make channels for sending and receiving our custom data types

	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16569, worldViewTx)
	go bcast.Receiver(16569, worldViewRx)

	// The example message. We just send one of these every second.
	// go func() {
	// 	helloMsg := WorldViewMsg{"Hello from " + id, 0}
	// 	for {
	// 		helloMsg.Iter++
	// 		helloTx <- helloMsg
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	// fmt.Println("Started")
	// for {
	// 	select {
	// 	case p := <-peerUpdateCh:
	// 		fmt.Printf("Peer update:\n")
	// 		fmt.Printf("  Peers:    %q\n", p.Peers)
	// 		fmt.Printf("  New:      %q\n", p.New)
	// 		fmt.Printf("  Lost:     %q\n", p.Lost)

	// 	case a := <-helloRx:
	// 		fmt.Printf("Received: %#v\n", a)
	// 	}
	// }
}
