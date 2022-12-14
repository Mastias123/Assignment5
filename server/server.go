package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"

	proto "github.com/Mastias123/Assignment5.git/grpc"
)

// go run server/server.go -sPort 5001
// go run server/server.go -sPort 5002
// go run server/server.go -sPort 5003
var wgMain sync.WaitGroup

type Server struct {
	proto.UnimplementedRegisterServer
	id          int32
	port        int
	maxBid      int
	maxBidId    int
	auctionOver bool
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

var sPort = flag.Int("sPort", 0, "server port number")
var bidders []bidder

func main() {
	flag.Parse()

	server1 := &Server{
		id:          1,
		port:        *sPort,
		maxBid:      0,
		maxBidId:    0,
		auctionOver: false,
	}

	go startServer(server1)

	wgMain.Add(1)
	wgMain.Wait()

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

func (s *Server) JoinServer(rq *proto.Request, rjss proto.Register_JoinServerServer) error {
	log.Printf("ID %d Connected to server id %d", rq.Id, s.id)
	var channel = make(chan bool, 1)
	cl := Client{int32(rq.Id), int32(rq.Port), rjss}
	cl.stream.Send(&proto.Reply{Id: cl.clientId, Msg: ""})

	<-channel //makes sure that join server isn't exited
	return nil
}

func (s *Server) PlaceBid(con context.Context, b *proto.Bid) (*proto.Conformation, error) {

	if s.auctionOver == true {
		return &proto.Conformation{Comment: "auction is Over ", MaxBid: int32(s.maxBid), MaxBidId: int32(s.maxBidId)}, nil
	}

	if b.MyPerseptionOfTheActonsMaxBid < int32(s.maxBid) { //If the bidder doesn't know what the current highest bid is

		return nil, errors.New("you do not know what the current max bid is")
	}

	log.Printf("Server -%d- resived bid from id %d", s.id, b.ClientId)
	b.Amount = b.Amount + int32(s.maxBid)
	log.Printf("Amount: %d", b.Amount)

	bidder := bidder{b.ClientId, b.Amount, b.ClientPort}

	if !contains(bidders, bidder.bidderId) { //Register the first time a bidder has placed a bid
		bidders = append(bidders, bidder)
	}
	opdateHighestBid(bidders, b.Amount, b.ClientId)

	s.maxBidId = int(b.ClientId)
	s.maxBid = int(b.Amount)
	if s.maxBid == 200 {
		s.auctionOver = true
	}
	log.Printf("max bid: %d", s.maxBid)
	return &proto.Conformation{Comment: "success", MaxBid: int32(s.maxBid), MaxBidId: int32(s.maxBidId)}, nil
}

func (s *Server) Result(con context.Context, rr *proto.ResultRequest) (*proto.Auctionresult, error) {

	if s.auctionOver {
		return &proto.Auctionresult{Id: int32(s.maxBidId), MaxBid: int32(s.maxBid), IsOver: s.auctionOver}, nil
	}

	return &proto.Auctionresult{Id: int32(s.maxBidId), MaxBid: int32(s.maxBid), IsOver: s.auctionOver}, nil
}

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
