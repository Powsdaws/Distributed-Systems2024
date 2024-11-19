package main

import (
	proto "AuctionSystem/grpc"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"encoding/json"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var id uint32

func main() {

	log.Println("Enter an id for the client")
	fmt.Scanln(&id)

	ConnectServers()
	CLI()

}

var connections []proto.AuctionServiceClient

func ConnectServers() {
	conn1, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client1 := proto.NewAuctionServiceClient(conn1)
	conn2, err := grpc.NewClient("localhost:5051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client2 := proto.NewAuctionServiceClient(conn2)
	conn3, err := grpc.NewClient("localhost:5052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client3 := proto.NewAuctionServiceClient(conn3)
	connections = append(connections, client1)
	connections = append(connections, client2)
	connections = append(connections, client3)
}

func sendBid(amount uint32) {
	log.Println(amount)
	bid := proto.Bid{
		ClientId: id,
		Amount:   amount,
	}

	var responses []*proto.Ack
	var errors []error

	//We make a bid to every node in the server
	for _, conn := range connections {
		ack, err := conn.Bid(context.Background(), &bid)
		if err != nil {
			errors = append(errors, err)
		}
		responses = append(responses, ack)
	}

	if len(errors) != 0 {
		var agreedError = getMostOccurring(errors)
		log.Println(agreedError.Error())
		return
	}

	var agreedResponse = getMostOccurring(responses)

	if agreedResponse.Acknowledgement {
		log.Println("Bid successful!    Amount: ", amount)
	} else {
		log.Println("Bid failed!")
	}
}

func sendResult() {

	var results []*proto.Outcome

	for _, conn := range connections {
		result, _ := conn.Result(context.Background(), &proto.Empty{})
		results = append(results, result)
	}

	var agreedResult = getMostOccurring(results)

	if agreedResult.WinnerId != 0 {
		log.Println("Winner: ", agreedResult.WinnerId, "!    Amount: ", agreedResult.HighestBid)
	} else {
		log.Println("No winner!    Current highest bid: " + strconv.Itoa(int(agreedResult.HighestBid)))
	}
}

func getMostOccurring[T any](elements []T) T {
	var hashes = make(map[string](T))
	var occurences = make(map[string](int))

	for _, element := range elements {
		// serialize and hash the result, increment occurences
		json, _ := json.Marshal(element)

		hashes[string(json)] = element
		occurences[string(json)] = occurences[string(json)] + 1
	}

	// get most occuring hash
	var mostOccuringHash string
	var mostOccuringOccurence int
	for hash, occurence := range occurences {
		if occurence > mostOccuringOccurence {
			mostOccuringHash = hash
			mostOccuringOccurence = occurence
		}
	}

	return hashes[mostOccuringHash]
}

func CLI() {

	reader := bufio.NewReader(os.Stdin)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSuffix(input, "\n")
		input = strings.TrimSuffix(input, "\r") //FOOOD >TIME>>>>>>>

		splitInput := strings.Split(input, " ")

		if splitInput[0] == "Bid" {

			var amount uint64

			amount, err = strconv.ParseUint(splitInput[1], 10, 32)
			if err != nil {
				panic(err)
			}
			log.Println("amount", amount)
			log.Println("Bid: ", splitInput[1])
			sendBid(uint32(amount))
		} else if splitInput[0] == "Result" {
			sendResult()
		} else {
			log.Println("Invalid command")
		}

	}
}
