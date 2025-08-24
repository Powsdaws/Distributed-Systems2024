package main

import "fmt"

func main() {
	var A, B, C float32
	fmt.Scan(&A, &B, &C)

	var I, J, K float32
	fmt.Scan(&I, &J, &K)

	var cocktails float32 = 0

	//Figure out how many cocktails that can be made:
	var cA float32 = A / I
	cB := B / J
	cC := C / K

	cocktails = min(cA, cB, cC)

	fmt.Println(A-(I*cocktails), B-(J*cocktails), C-(K*cocktails))

}
