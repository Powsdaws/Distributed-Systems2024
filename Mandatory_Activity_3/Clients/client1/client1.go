package main

import (
	proto "Chitty-Chat/grpc"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Not working")
	}

	client := proto.NewChatServiceClient(conn)

	var currentLambort uint32 = 0
	lambortTime := &proto.Timestamp{
		LamportTime: currentLambort,
	}

	connection, err := client.Subscribe(context.Background(), lambortTime)
	if err != nil {
		log.Fatalf("Connection not established")
	}

	//channels for communcating the current lambortTime
	channelIn := make(chan uint32)
	channelOut := make(chan uint32)

	//Thread for writing messages in the cli
	go sendMessages(client, &channelIn, &channelOut)

	//This recieves messages and listens
	for {

		//Keeps track of lambort if this clients sends a message
		select {
		case newLambort := <-channelOut:
			if newLambort > currentLambort {
				currentLambort = newLambort
			}
		case <-time.After(2 * time.Millisecond):
			fmt.Println("no new lambort time from server")
		}

		//We check if we received a new broadcast
		var recMessage, err = connection.Recv()

		if err != nil {
			log.Fatalf("recv wrong")
		}

		if recMessage != nil {
			fmt.Println("client1 received: ", recMessage.Text)
			fmt.Println("client1 received time: ", recMessage.LamportTime)
			if recMessage.LamportTime > lambortTime.LamportTime {
				lambortTime.LamportTime = recMessage.LamportTime + 1
			} else {
				lambortTime.LamportTime += 1
			}
			fmt.Println("client1 lambort: ", lambortTime.LamportTime)
		}

	}

}

func sendMessages(client proto.ChatServiceClient, channelIn *chan uint32, channelOut *chan uint32) {
	var currentLambort uint32 = 0
	for {
		var newMessage = ""
		fmt.Scanln(&newMessage)

		currentLambort = currentLambort + 1

		select {
		case newLambort := <-*channelIn:
			if newLambort > currentLambort {
				currentLambort = newLambort
			}
		case <-time.After(2 * time.Millisecond):
			fmt.Println("no new lambort time from client itself")
		}

		//TODO: Fix this so its does not get stuck : ) - we need to make sure the main thread knows that we
		//updated the lamport time in this thread
		//*channelOut <- currentLambort

		newMessageChat := &proto.ChatMessage{
			Text:        newMessage,
			LamportTime: uint32(currentLambort),
		}

		client.Publish(context.Background(), newMessageChat)

	}
}
