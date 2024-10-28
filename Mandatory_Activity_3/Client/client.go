package main

import (
	proto "Chitty-Chat/grpc"
	"bufio"
	"context"
	"log"
	"os"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var lamportLock sync.Mutex
var lamportTime uint32

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Not working")
	}

	client := proto.NewChatServiceClient(conn)

	// Subscribe to messages from Server
	lamportLock.Lock()
	lamportTime += 1
	//log.Println("participant sent subscribe request with time: " + strconv.Itoa(int(lamportTime)))
	connection, err := client.Subscribe(context.Background(), &proto.Timestamp{LamportTime: lamportTime})
	lamportLock.Unlock()
	if err != nil {
		log.Fatalf("Connection not established")
	}

	//Thread for writing messages in the cli
	go sendMessages(client)

	//This recieves messages and listens
	for {
		//We check if we received a new broadcast
		var recMessage, err = connection.Recv()

		if err != nil {
			log.Fatalf("Server closed unexpectedly")
		}

		if recMessage != nil {
			log.Println(recMessage.Text)

			syncTime(recMessage.LamportTime)

		}

	}

}

func syncTime(recvTime uint32) {
	lamportLock.Lock()
	lamportTime = max(lamportTime, recvTime) + 1
	//log.Println("Client time: ", lamportTime)
	lamportLock.Unlock()
}

func sendMessages(client proto.ChatServiceClient) {

	reader := bufio.NewReader(os.Stdin)

	for {
		newMessage, err := reader.ReadString('\n')

		if err != nil {
			log.Println("Unsubscribed") //we only get this when we unsubscribe
			continue
		}
		if len(newMessage) > 128 {
			log.Println("Message too long :(")
			continue
		}

		//We lock lamport field
		lamportLock.Lock()
		lamportTime = lamportTime + 1

		newMessageChat := &proto.ChatMessage{
			Text:        newMessage,
			LamportTime: lamportTime,
		}

		//log.Println("client sent message: \"", newMessage, "\" with time ", lamportTime)

		client.Publish(context.Background(), newMessageChat)
		lamportLock.Unlock()
	}
}
