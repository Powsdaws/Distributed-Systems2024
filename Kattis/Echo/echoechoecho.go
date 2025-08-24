package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin) //Creates a new scanner that can read text

	if scanner.Scan() {
		input := scanner.Text()
		for i := 0; i < 3; i++ {
			fmt.Print(input, " ")
		}
	}

	//for i := 0; i < len(arg); i++ {
	//	fmt.Println(i, arg[i])
	//}
}
