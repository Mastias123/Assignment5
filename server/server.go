package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	proto "github.com/Mastias123/Assignment5.git/grpc"
)

type Server struct {
	proto.UnsafeRegisterServer
	id        int32
	clients   map[int32]proto.RegisterClient
	timestamp int32
}

func main() {

	server := &Server{
		id: 8,
	}

	go startServer(server)


}

func startServer(server *Server) {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", 5000))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v\n", err)
	}

	// Register the grpc server and serve its listener
	proto.RegisterRegisterServer(grpcServer, server)
	serveError := grpcServer.Serve(list)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}

	for i := 0; i < 3; i++ {
		port := int32(5000) + int32(i)

		if port == 5000 {
			continue
		}

		var conn *grpc.ClientConn
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s\n", err)
		}
		defer conn.Close()
		c := proto.NewRegisterClient(conn)
		server.clients[port] = c
	}
}
