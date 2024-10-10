package main

import (
	proto "Chitty-Chat/grpc"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Not working")
	}

	client := proto.NewChatServiceClient(conn)

	lambortTime := &proto.Timestamp{
		LamportTime: 0,
	}

	connection, err := client.Subscribe(context.Background(), lambortTime)
	if err != nil {
		log.Fatalf("Connection not established")
	}

	for {
		var recMessage, err = connection.Recv()

		if err != nil {
			log.Fatalf("recv wrong")
		}

		if recMessage != nil {
			fmt.Println("client1 received: ", recMessage.Text)
			if recMessage.LamportTime > lambortTime.LamportTime {
				lambortTime.LamportTime = recMessage.LamportTime + 1
			} else {
				lambortTime.LamportTime += 1
			}
			fmt.Println("client1 lambort: ", lambortTime.LamportTime)
		}

	}

}
