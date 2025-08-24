package main

import (
	"fmt"
)

func main() {
	var ok bool = true

	var input string
	fmt.Scan(&input)

	var m = map[string]int{}
	suits := []int{13, 13, 13, 13} //P, K, H, T

	//Map each input to a value for how many times it appeared.
	for i := 0; i < len(input); i = i + 3 {
		m[input[i:i+3]] = m[input[i:i+3]] + 1

		//Checks if too many of same card
		if m[input[i:i+3]] > 1 {
			fmt.Println("GRESKA")
			ok = false
			break
		} else {
			//Remove from remaining cards
			var suit_type string = input[i : i+1]

			switch suit_type {
			case "P":
				suits[0] = suits[0] - 1
			case "K":
				suits[1] = suits[1] - 1
			case "H":
				suits[2] = suits[2] - 1
			case "T":
				suits[3] = suits[3] - 1
			}

		}
	}

	if ok {
		for suit := range suits {
			fmt.Print(suits[suit], " ")
		}

	}

}
