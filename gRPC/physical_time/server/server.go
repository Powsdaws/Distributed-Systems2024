package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "TimeServer/grpc"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTimeServiceServer
}

func (s *server) GetTime(ctx context.Context, req *pb.TimeRequest) (*pb.TimeResponse, error) {
	return &pb.TimeResponse{CurrentTime: time.Now().Format(time.RFC3339Nano)}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTimeServiceServer(grpcServer, &server{})
	log.Println("TimeService running on :50051")
	grpcServer.Serve(lis)
}
