package main

import (
	"context"

	proto "Chitty-Chat/grpc"
)

// The fields that the server has
type ChatServiceServer struct {
	proto.UnimplementedChatServiceServer
	lambortTime int32
	messages    []string
	streamsIn   []proto.ChatService_ConnectClient
	streamsOut  []proto.ChatService_ConnectServer
}

// implement Chat_message methods to be called from client
// Should log the message and time and then broadcast to all clients
func (s *ChatServiceServer) Connect(ctx context.Context, streamIn *proto.ChatService_ConnectClient, streamOut *proto.ChatService_ConnectServer) (*proto.ChatService_ConnectServer, error) {
	return proto.ChatService_ConnectServer, nil
}
