package main

import (
	"fmt"
	"math/rand"
	"time"
)

// int arrays with format: [sequence no, ack. no]
var s2c = make(chan [2]int) // channel for communicating from server to client
var c2s = make(chan [2]int) // ^^ vice versa

func main() {
	//Thread acting as client
	go client()

	//Thread acting as server
	go server()

	//Thread actin like a lil bitch
	//go lil-bitch();
	time.Sleep(5 * time.Second)
}

func client() {
	seq := rand.Intn(100)
	fmt.Printf("Client seq is: %d\n", seq)
	//We do the handshake
	if !clientHandshake(seq) {
		//Handshake went bad... terminate?

	}

	println("Client: Handshake complete, ready to send data")
	//We send data if handshake OK...

}

func clientHandshake(seq int) bool {
	// send SYN
	c2s <- [2]int{seq, 0} // send seq, 0 to server
	fmt.Printf("Client sent SYN: [%d, %d]\n", seq, 0)
	time.Sleep(500 * time.Millisecond)

	// receive SYN-ACK
	message := <-s2c
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("Client received SYN-ACK: [%d, %d]\n", message[0], message[1])
	nextSeq := seq + 1

	if message[1] != nextSeq { // ensure message ack matches expected seq
		// error handle
		return false
	}

	// send ACK
	fmt.Printf("Client sent ACK: [%d, %d]\n", nextSeq, message[0]+1)
	c2s <- [2]int{nextSeq, message[0] + 1}
	time.Sleep(500 * time.Millisecond)

	return true

}

// SERVER SHITE
// SERVER SHITE
// SERVER SHITE

func server() {
	seq := rand.Intn(100)
	fmt.Printf("Server seq is %d\n", seq)
	for { // keeps looking for a successfhandshake

		request := <-c2s
		fmt.Printf("Server recieved SYN: [%d, %d] \n", request[0], request[1])
		success, _, _ := serverHandshake(seq, request[0])
		if success {
			break
		}
	}

}

func serverHandshake(seq int, ack int) (bool, int, int) {

	// send SYN-ACK
	s2c <- [2]int{seq, ack + 1}
	fmt.Printf("Server sent SYN-ACK: [%d, %d]\n", seq, ack+1)
	time.Sleep(500 * time.Millisecond)

	// recieve ACK
	message := <-c2s
	fmt.Printf("Server received ACK: [%d, %d]\n", message[0], message[1])
	nextSeq := seq + 1
	time.Sleep(500 * time.Millisecond)

	if message[1] != nextSeq { // ensure message ack matches expected seq
		// error handle
		return false, -1, -1
	}
	return true, nextSeq, message[0] + 1
}
