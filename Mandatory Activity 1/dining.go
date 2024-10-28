package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var forkChannels = [10]chan int{make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int)}
var filoChannels = [10]chan int{make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int)}
var wg sync.WaitGroup
var delay time.Duration

/*
EXPLANATION
We use different channels for receiving and sending information between philosophers and forks. Forks only write into forkchannels, and only read from filochannels, and vice versa for filomen.
This ensures no deadlock can occur, as there is no situation where two processes are reading or writing to the same channel. This functions like a lock.
We use a waitgroup to end the process, but it does not allow communication between philosophers or forks.
Who gets to eat is determined through a lottery process, though simple allowing them to pick up and drop forks after a short delay would achieve the same.
The lottery does reduce the risk of everyone getting only one fork, even though this does not lead to a deadlock in our system.
*/

func main() {
	// get delay in ms from program arguments
	if len(os.Args) < 2 {
		delay = 0
	} else {
		delayArg, err := strconv.Atoi(os.Args[1])
		if err != nil {
			delayArg = 0
		}
		delay = time.Duration(delayArg) * time.Second
	}

	// we use wait groups to terminate the program once everyone has eaten 3 times
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go fork(i)
		go filoMan(i)
	}

	wg.Wait()
	fmt.Print("\n Everyone is fat now")

}

func fork(id int) { //
	var forkchannelLeft = forkChannels[((2*id)+1)%10]  //The channel used send information to left filoman
	var forkchannelRight = forkChannels[((2*id)+2)%10] //The channel used send information to right filoman
	var filoChannelLeft = filoChannels[((2*id)+1)%10]  //The channel used get information from left filoman
	var filoChannelRight = filoChannels[((2*id)+2)%10] //The channel used get information from right filoman
	for {
		messageLeft := <-filoChannelLeft   //Listen for message on left channel from left filoMan
		messageRight := <-filoChannelRight //Listen for message on right channel from right filoMan

		if messageLeft > messageRight { //Checks which filoman has the greatest priority, sending back 1 if that filoman should have the fork.
			forkchannelLeft <- 1
			forkchannelRight <- 0
		} else {
			forkchannelLeft <- 0
			forkchannelRight <- 1
		}

	}

}

func filoMan(id int) {
	var stomach int
	var filochannelLeft = filoChannels[((2 * id) + 1)] //The channel used send information to left fork
	var filochannelRight = filoChannels[2*id]          //The channel used send information to right fork
	var forkchannelLeft = forkChannels[((2 * id) + 1)] //The channel used get information from left fork
	var forkchannelRight = forkChannels[2*id]          //The channel used get information from right fork

	var priority int

	for {
		priority = rand.Intn(math.MaxInt) // generates a random number to be used in the fork lotteries
		time.Sleep(delay)
		// try to grab both forks
		filochannelLeft <- priority  //Sends priority via filoChannel to left fork
		filochannelRight <- priority //Sends priority via filoChannel to right fork

		answerLeft := <-forkchannelLeft   //Retrives answer from the left fork via left fork channel
		answerRight := <-forkchannelRight //Retrives answer from the right fork via right fork channel

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
