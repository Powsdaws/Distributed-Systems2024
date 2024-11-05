package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"strings"
	"time"

	proto "modfile/grpc"

	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Tokenreciever struct {
	proto.UnimplementedTokenServiceServer
}

var hasToken bool
var wantsToken bool

// this "happens" when another client sends the token to this client
func (s *Tokenreciever) SendToken(ctx context.Context, in *proto.Token) (*proto.Empty, error) {
	log.Println("This client recieved the token")
	hasToken = true
	return &proto.Empty{}, nil
}

func main() {

	reciever := Tokenreciever{}
	//set up own node

	log.Println("Declare a host address for this node")
	var hostAddress string

	fmt.Scanln(&hostAddress)

	//hostAddress, err := reader.ReadString("")
	var hostAddressString = ":" + hostAddress
	//hostAddressString = strings.Trim(hostAddressString, "\n")

	go reciever.startListener(&hostAddressString)

	//set up the address for the next node in line
	log.Println("Declare a recieving address")
	var recieverAddress string
	fmt.Scanln(&recieverAddress)

	recieverAddressString := strings.Trim("localhost:"+recieverAddress, "\n")
	conn, err := grpc.NewClient(recieverAddressString, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Not working")
	}

	sender := proto.NewTokenServiceClient(conn)

	//the client with the adress 5050 always starts with the token
	if hostAddressString == ":5052" {
		hasToken = true
		log.Println("This client has started with the token")
	}

	for {
		if hasToken {
			time.Sleep(1 * time.Second)
			if wantsToken {
				log.Println("TOKEN USED")
				time.Sleep(1 * time.Second)
				wantsToken = false
			} else {
				log.Println("Did not want to access critical section")
			}

			_, err := sender.SendToken(context.Background(), &proto.Token{})
			if err != nil {
				log.Println(err)
			}
			log.Println("Token has been forwarded")
			hasToken = false
		}

		if !wantsToken {
			desire := rand.IntN(10)
			if desire < 2 {
				wantsToken = true
				log.Println("Client now wants the token")
			}
			time.Sleep(1 * time.Second)
		}
		time.Sleep(1 * time.Second)
	}

}

func (s *Tokenreciever) startListener(address *string) {
	grpcReciever := grpc.NewServer()

	listener, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}

	proto.RegisterTokenServiceServer(grpcReciever, s)

	error := grpcReciever.Serve(listener)
	if error != nil {
		log.Fatalf("Failed")
	}

}
