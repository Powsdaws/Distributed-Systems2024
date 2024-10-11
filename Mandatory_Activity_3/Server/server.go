package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	proto "Chitty-Chat/grpc"

	"google.golang.org/grpc"
)

// The fields that the server has
type ChatServiceServer struct {
	proto.UnimplementedChatServiceServer
	lamportTime   uint32
	newestMessage string
	subscriptions []proto.ChatService_SubscribeServer
	subcount      uint32
}

func (s *ChatServiceServer) Publish(ctx context.Context, msg *proto.ChatMessage) (*proto.Empty, error) {
	s.lamportTime += 1
	fmt.Println("Server received message: ", msg.Text)
	s.Broadcast(msg)

	return &proto.Empty{}, nil
}

// implement Chat_message methods to be called from client
// Should log the message and time and then broadcast to all clients
func (s *ChatServiceServer) Subscribe(timestamp *proto.Timestamp, stream proto.ChatService_SubscribeServer) error {
	s.subscriptions = append(s.subscriptions, stream)
	s.syncTime(timestamp.LamportTime) // ensure the clock is synced

	s.subcount = s.subcount + 1
	clientID := s.subcount // print this number when the client joins and when they leave

	fmt.Println("A client has joined")

	joinMessage := &proto.ChatMessage{
		Text:        "Client " + strconv.Itoa(int(clientID)) + " has subscribed",
		LamportTime: s.lamportTime,
	}
	fmt.Println("text: ", joinMessage.Text)
	s.Broadcast(joinMessage)

	<-stream.Context().Done()
	//make sure to do whatever requirements say when unsubscribed

	unsubscribeMessage := &proto.ChatMessage{
		Text:        ("Client " + strconv.Itoa(int(clientID)) + " has unsubcribed"),
		LamportTime: s.lamportTime,
	}
	s.Broadcast(unsubscribeMessage)

	return nil
}

func (s *ChatServiceServer) Broadcast(message *proto.ChatMessage) {
	for _, sub := range s.subscriptions {
		sub.Send(message)
	}

}

func (s *ChatServiceServer) syncTime(recvTime uint32) {
	s.lamportTime = max(s.lamportTime, recvTime) + 1
	fmt.Println("Server lambort: ", s.lamportTime)
}

func main() {
	server := ChatServiceServer{}
	server.start_server()
}

func (s *ChatServiceServer) start_server() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Did not work")
	}

	proto.RegisterChatServiceServer(grpcServer, s)

	fmt.Println("Server started")
	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("Did not work")
	}

}
