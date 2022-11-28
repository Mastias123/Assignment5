package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	proto "github.com/Mastias123/Assignment5.git/grpc"
)

type Server struct {
	proto.UnimplementedRegisterServer
	id int32
	//bidders   []Client
	timestamp int32
	port      int
	maxBid    int
	maxBidId  int
	crash     bool
}
type Client struct {
	clientId   int32
	clientPort int32
	stream     proto.Register_JoinServerServer
}

type bidder struct {
	bidderId   int32
	highestBid int32
	bidderPort int32
}

var bidders []bidder

func main() {

	server1 := &Server{
		id:        1,
		timestamp: 0,
		port:      5001,
		maxBid:    0,
		maxBidId:  0,
		crash:     false,
	}
	server2 := &Server{
		id:        2,
		timestamp: 0,
		port:      5002,
		maxBid:    0,
		maxBidId:  0,
		crash:     true,
	}

	server3 := &Server{
		id:        3,
		timestamp: 0,
		port:      5003,
		maxBid:    0,
		maxBidId:  0,
		crash:     false,
	}

	go startServer(server1)
	go startServer(server2)
	go startServer(server3)

	for {

	}

}

func startServer(server *Server) {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", server.port))

	if err != nil {
		log.Fatalf("Failed to listen on port: %v\n", err)
	}

	log.SetFlags(0)
	log.Printf("Started server at port: %d\n", server.port)

	// Register the grpc server and serve its listener
	proto.RegisterRegisterServer(grpcServer, server)

	serveError := grpcServer.Serve(list)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}

// The join server function is named after the grpc function, and when you run the proto command the proto file will create a function signature that has to be implemented
func (s *Server) JoinServer(rq *proto.Request, rjss proto.Register_JoinServerServer) error {
	log.Printf("ID %d Connected to server id %d", rq.Id, s.id)
	var channel = make(chan bool, 1)
	cl := Client{int32(rq.Id), int32(rq.Port), rjss}
	cl.stream.Send(&proto.Reply{Id: cl.clientId, Msg: ""})

	<-channel //makes sure that join server isn't exited
	return nil
}

func (s *Server) PlaceBid(con context.Context, b *proto.Bid) (*proto.Conformation, error) {

	if s.crash == true {

	}

	if b.MyPerseptionOfTheActonsMaxBid < int32(s.maxBid) { //If the bidder doesn't know what the current highest bid is

		return nil, errors.New("you do not know what the current max bid is")
	}

	log.Printf("Server -%d- resived bid from id %d", s.id, b.ClientId)
	b.Amount = b.Amount + int32(s.maxBid)
	log.Printf("Amount: %d", b.Amount)

	bidder := bidder{b.ClientId, b.Amount, b.ClientPort}

	if !contains(bidders, bidder.bidderId) { //Register the first time a bidder had placed a bid
		bidders = append(bidders, bidder)
	}
	opdateHighestBid(bidders, b.Amount, b.ClientId)

	s.maxBidId = int(b.ClientId)
	s.maxBid = int(b.Amount)
	log.Printf("max bid: %d", s.maxBid)
	return &proto.Conformation{Comment: "success ", MaxBid: int32(s.maxBid), MaxBidId: int32(s.maxBidId)}, nil
}

func (s *Server) Result(con context.Context, rr *proto.ResultRequest) (*proto.Auctionresult, error) {

	return &proto.Auctionresult{Id: int32(s.maxBidId), MaxBid: int32(s.maxBid)}, nil
}

//_____________________________________________________________
//_____________________________________________________________
//_____________________________________________________________

func contains(b []bidder, bId int32) bool {
	for _, v := range b {
		if v.bidderId == bId {
			return true
		}
	}
	return false
}

func opdateHighestBid(b []bidder, hBid int32, bId int32) {
	for _, v := range b {
		if v.bidderId == bId {
			if v.highestBid < hBid {
				v.highestBid = hBid
			}
		}
	}
}
