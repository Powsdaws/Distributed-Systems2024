package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

var forks = [10]chan int{make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int), make(chan int)}

var wg sync.WaitGroup

func main() {

	wg.Add(5)
	for i := 0; i < 5; i++ {
		go fork(i)
		go filoMan(i)
	}

	wg.Wait()
	fmt.Print("\n Everyone is fat now")

}

func fork(id int) {
	var forkLeft = forks[((2*id)+1)%10]  //The left channel used to communicate with the left filoman
	var forkRight = forks[((2*id)+2)%10] //The right channel used to communicate with the right filo
	for {
		messageLeft := <-forkLeft
		messageRight := <-forkRight

		if messageLeft > messageRight {
			forkLeft <- 1
			forkRight <- 0
		} else {
			forkLeft <- 0
			forkRight <- 1
		}
		time.Sleep(200)
	}

}

func filoMan(id int) {
	var stomach int
	var forkLeft = forks[((2 * id) + 1)] //The left channel used to communicate with the left fork
	var forkRight = forks[2*id]          //The right channel used to communicate with the right fork

	var priority int

	for { // while hungry
		if stomach < 3 {
			priority = rand.Intn(math.MaxInt)
		} else {
			priority = 0
		}

		// try to eat
		forkLeft <- priority
		forkRight <- priority

		time.Sleep(100)
		answerLeft := <-forkLeft
		answerRight := <-forkRight

		if answerLeft == 1 && answerRight == 1 {
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
