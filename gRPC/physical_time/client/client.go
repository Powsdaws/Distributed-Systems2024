package main

import (
	"context"
	"log"
	"time"

	pb "TimeServer/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTimeServiceClient(conn)
	start := time.Now()
	resp, err := client.GetTime(context.Background(), &pb.TimeRequest{})
	if err != nil {
		log.Fatalf("Error calling GetTime: %v", err)
	}
	rtt := time.Since(start)

	log.Printf("Received time: %s | RTT: %s", resp.CurrentTime, rtt)
}
