package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

var channels = [10]chan int{make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int)}

var wg sync.WaitGroup

func main() {
	// we use wait groups to terminate the program once everyone has eaten 3
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go fork(i)
		go filoMan(i)
	}

	wg.Wait()
	fmt.Print("\n Everyone is fat now")

}

func fork(id int) { //
	var channelLeft = channels[((2*id)+1)%10]  //The left channel used to communicate with left filoman
	var channelRight = channels[((2*id)+2)%10] //The right channel used to communicate with left filoman
	for {
		messageLeft := <-channelLeft
		messageRight := <-channelRight

		if messageLeft > messageRight {
			channelLeft <- 1
			channelRight <- 0
		} else {
			channelLeft <- 0
			channelRight <- 1
		}
		time.Sleep(200)
	}

}

func filoMan(id int) {
	var stomach int
	var channelLeft = channels[((2 * id) + 1)] //The left channel used to communicate with the left fork
	var channelRight = channels[2*id]          //The right channel used to communicate with the right fork

	var priority int

	for {
		priority = rand.Intn(math.MaxInt) // generates a random number to be used in the fork lotteries

		// try to grab both forks
		channelLeft <- priority
		channelRight <- priority

		time.Sleep(100) // wait for answers to arrive
		answerLeft := <-channelLeft
		answerRight := <-channelRight

		if answerLeft == 1 && answerRight == 1 { //if both forks are available, eat
			stomach++
			fmt.Printf("Filoman %d grasps both forks forcefully\n", id)
			fmt.Printf("Filoman %d ate. Stomach at %d/3\n", id, stomach)
			if stomach == 3 {
				wg.Done()
			}
		} else {
			fmt.Printf("Filoman %d is thinking\n", id)
		}

	}

}
