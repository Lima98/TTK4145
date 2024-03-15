package network

import (
	"Elevator_project/network/network/bcast"
	"Elevator_project/network/network/peers"
	wv "Elevator_project/worldviewmessage"

)




func Network(worldViewTx chan wv.WorldViewMsg, worldViewRx chan wv.WorldViewMsg, peerUpdateCh chan peers.PeerUpdate, id string) {

	peerTxEnable := make(chan bool)

	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	go bcast.Transmitter(16569, worldViewTx)
	go bcast.Receiver(16569, worldViewRx)

}
