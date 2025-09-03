package main

import "time"

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go pingerA(ch1, ch2)
	go pingerB(ch1, ch2)

	for { //Loop so the program does not terminate before the goroutines start
	}

}

func pingerA(ch1 chan string, ch2 chan string) {
	for {
		ch1 <- "ping"
		time.Sleep(100 * time.Millisecond)
		read := <-ch2
		println("A reviewed: ", read)
	}

}

func pingerB(ch1 chan string, ch2 chan string) {
	for {
		read := <-ch1
		println("B reviewed: ", read)
		time.Sleep(100 * time.Millisecond)
		ch2 <- "pong"
	}
}
