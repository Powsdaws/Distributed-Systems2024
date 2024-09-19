package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// int arrays with format: [sequence no, ack. no]
var s2m = make(chan [2]int) // channel for communicating from server to client
var c2m = make(chan [2]int) // ^^ vice versa

var m2c = make(chan [2]int)
var m2s = make(chan [2]int)

const dropPackageChance = 25 // the chance to drop the package in percent

var wg sync.WaitGroup

func main() {

	wg.Add(2)

	//Thread acting as client
	go client()

	//Thread acting as server
	go server()

	//Thread for the middle man (can drop packages)
	go middleMan()

	wg.Wait()
}

func client() {
	defer wg.Done()

	seq := rand.Intn(100)
	fmt.Printf("Client seq is: %d\n", seq)
	//We do the handshake
	if !clientHandshake(seq) {
		//Handshake went bad... terminate?
		println("Handshake FAILED")
	}

	println("Client: Handshake complete, ready to send data")
	//We send data if handshake OK...
}

func clientHandshake(seq int) bool {
	// send SYN
step1:

	var pack = [2]int{seq, 0}
	fmt.Printf("Client sent SYN: [%d, %d]\n", seq, 0)
	sendPackage(&c2m, pack)
	//time.Sleep(500 * time.Millisecond)

	// receive SYN-ACK
	var message [2]int
	select {
	case message = <-m2c:
		goto step2
	case <-time.After(1 * time.Second):
		fmt.Println("Client: SYN-ACK package never recieved... resending SYN")
		goto step1
	}

step2:

	//time.Sleep(500 * time.Millisecond)
	fmt.Printf("Client received SYN-ACK: [%d, %d]\n", message[0], message[1])
	nextSeq := seq + 1

	if message[1] != nextSeq { // ensure message ack matches expected seq
		// error handle
		return false
	}

	// send ACK
	fmt.Printf("Client sent ACK: [%d, %d]\n", nextSeq, message[0]+1)
	sendPackage(&c2m, [2]int{nextSeq, message[0] + 1})
	//time.Sleep(500 * time.Millisecond)

	select {
	case message = <-m2c:
		goto step2
	case <-time.After(3 * time.Second):
		return true
	}

}

func sendPackage(ch *chan [2]int, pack [2]int) {
	*ch <- pack
}

// SERVER SHITE
// SERVER SHITE
// SERVER SHITE

func server() {
	defer wg.Done()

	seq := rand.Intn(100)
	fmt.Printf("Server seq is %d\n", seq)

	for {
		lastRequest := -1
		// keep waiting for SYN
		request := <-m2s

		// if not SYN, ignore
		if request[1] != 0 {
			continue
		}

		fmt.Printf("Server recieved SYN: [%d, %d] \n", request[0], request[1])

		if request[0] == lastRequest { // checks if the current request is the same as last recieved request
			fmt.Printf("Server: request already being processed... ignoring... \n")
			continue
		}
		lastRequest = request[0]
		success, _, _ := serverHandshake(seq, request[0])
		if success {
			break
		}
	}
}

func serverHandshake(seq int, ack int) (bool, int, int) {

	// send SYN-ACK
step1:
	fmt.Printf("Server sent SYN-ACK: [%d, %d]\n", seq, ack+1)
	sendPackage(&s2m, [2]int{seq, ack + 1})
	//time.Sleep(50 * time.Millisecond)

	// recieve ACK
	var message [2]int
	select {
	case message = <-m2s:
		goto step2
	case <-time.After(1 * time.Second):
		fmt.Println("Server: ACK package never recieved... resending SYN-ACK")
		goto step1
	}

step2:
	fmt.Printf("Server received ACK: [%d, %d]\n", message[0], message[1])
	nextSeq := seq + 1
	//time.Sleep(50 * time.Millisecond)

	if message[1] != nextSeq { // ensure message ack matches expected seq
		// error handle
		return false, -1, -1
	}
	return true, nextSeq, message[0] + 1
}

func middleMan() {
	for {
		select {
		//wating for message from server
		case msg1 := <-s2m:
			if rand.Intn(100) < dropPackageChance {
				fmt.Println("Oops.. a package was lost")
				continue
			}
			m2c <- msg1
			//wating for message from client
		case msg2 := <-c2m:
			if rand.Intn(100) < dropPackageChance {
				fmt.Println("Oops.. a package was lost")
				continue
			}
			m2s <- msg2
		}
	}
}
