package main

import (
	proto "AuctionSystem/grpc"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type AuctionServiceServer struct {
	proto.UnimplementedAuctionServiceServer
	EndAt         time.Time
	HighestBid    uint32
	HighestBidder uint32
}

func (s *AuctionServiceServer) Bid(ctx context.Context, bid *proto.Bid) (*proto.Ack, error) {

	//Start Auction if not prior bid has been placed
	if s.EndAt.IsZero() {
		s.startAuction()
	}

	var errorMessage error

	//We check if auction is over
	currentTime := time.Now()
	if currentTime.After(s.EndAt) {
		ack := &proto.Ack{
			Acknowledgement: false,
		}
		errorMessage = errors.New("Bid failed: Auction is over")
		return ack, errorMessage
	}

	//Check if bid is higher than the current highest bid
	if bid.Amount > s.HighestBid {
		s.HighestBid = bid.Amount
		s.HighestBidder = bid.ClientId
	} else {
		errorMessage = errors.New("Bid failed: Bid not high enough")
		ack := &proto.Ack{
			Acknowledgement: false,
		}
		return ack, errorMessage
	}

	ack := &proto.Ack{
		Acknowledgement: true,
	}

	return ack, nil
}

func (s *AuctionServiceServer) Result(ctx context.Context, empty *proto.Empty) (*proto.Outcome, error) {
	outcome := &proto.Outcome{
		HighestBid: s.HighestBid,
	}

	// if auction is over, add the winner else add the current highest bid
	outcome.WinnerId = s.HighestBidder

	return outcome, nil
}

func main() {
	// get a host port from the user

	log.Println("Declare a host port for this server, e.g.: 5050")
	var port string
	fmt.Scanln(&port)

	// start an auction
	server := AuctionServiceServer{}
	server.startServer(port)

}

func (s *AuctionServiceServer) startServer(port string) {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Did not work")
	}

	proto.RegisterAuctionServiceServer(grpcServer, s)

	log.Println("A auction server/node has startd")

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("Something failed while starting the server node")
	}
}

func (s *AuctionServiceServer) startAuction() {
	s.EndAt = time.Now().Add(time.Second * 30)
}
