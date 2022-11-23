package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	proto "github.com/Mastias123/Assignment5.git/grpc"
)

type Server struct {
	proto.UnimplementedRegisterServer
	id        int32
	clients   map[int32]proto.RegisterClient
	timestamp int32
	port      int
}

func main() {

	server1 := &Server{
		id:        1,
		timestamp: 0,
		port:      5001,
	}
	server2 := &Server{
		id:        2,
		timestamp: 0,
		port:      5002,
	}

	server3 := &Server{
		id:        3,
		timestamp: 0,
		port:      5003,
	}

	//If you want to run the function as a go routine you have to make sure that this main function does not terminate. This can be done by eather creating a infinite forloop og a wait group
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

	for {
	}
	return nil
}
