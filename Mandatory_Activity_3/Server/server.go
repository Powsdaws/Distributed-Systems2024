package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"sync"

	proto "Chitty-Chat/grpc"

	"google.golang.org/grpc"
)

// The fields that the server has
type ChatServiceServer struct {
	proto.UnimplementedChatServiceServer
	lamportTime   uint32
	subscriptions []proto.ChatService_SubscribeServer
	subcount      uint32
	lock          sync.Mutex
}

func (s *ChatServiceServer) Publish(ctx context.Context, msg *proto.ChatMessage) (*proto.Empty, error) {
	s.lock.Lock()
	log.Println("Server received message: " + msg.Text + " with time: " + strconv.Itoa(int(msg.LamportTime)))
	s.syncTime(msg.LamportTime)

	s.lamportTime += 1 //increment time before sending message
	msg.LamportTime = s.lamportTime
	s.Broadcast(msg)
	s.lock.Unlock()
	return &proto.Empty{}, nil
}

// implement Chat_message methods to be called from client
// Should log the message and time and then broadcast to all clients
func (s *ChatServiceServer) Subscribe(timestamp *proto.Timestamp, stream proto.ChatService_SubscribeServer) error {
	s.lock.Lock()
	s.subscriptions = append(s.subscriptions, stream)

	s.syncTime(timestamp.LamportTime) // ensure the clock is synced
	s.subcount = s.subcount + 1
	clientID := s.subcount // print this number when the client joins and when they leave
	text := "Participant " + strconv.Itoa(int(clientID)) + " joined Chitty-Chat at Lamport Time " + strconv.Itoa(int(s.lamportTime))
	log.Println(text)

	s.lamportTime += 1 // increment time cause we are about to send a message

	joinMessage := &proto.ChatMessage{
		Text:        text,
		LamportTime: s.lamportTime,
	}

	s.Broadcast(joinMessage)
	s.lock.Unlock()

	//If client closes stream disconnect
	<-stream.Context().Done()
	//make sure to do whatever requirements say when unsubscribed

	s.lock.Lock()
	s.lamportTime++

	unsubscribeMessage := &proto.ChatMessage{
		Text:        ("Participant " + strconv.Itoa(int(clientID)) + " left Chitty-Chat at Lamport Time " + strconv.Itoa(int(s.lamportTime))),
		LamportTime: s.lamportTime,
	}
	log.Println("Server recieved unsubscribe request with time: " + strconv.Itoa(int(s.lamportTime)))
	s.Broadcast(unsubscribeMessage)

	s.lock.Unlock()

	return nil
}

func (s *ChatServiceServer) Broadcast(message *proto.ChatMessage) {
	//log.Println("Broadcasting: ", message.Text,)
	for _, sub := range s.subscriptions {
		sub.Send(message)
	}
}

func (s *ChatServiceServer) syncTime(recvTime uint32) {
	s.lamportTime = (max(s.lamportTime, recvTime)) + 1
	//log.Println("Server lamport synced to: ", s.lamportTime)
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

	log.Println("Server started")
	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("Did not work")
	}

}
