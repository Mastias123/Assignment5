package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"strconv"

	proto "github.com/Mastias123/Assignment5.git/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	proto.UnimplementedRegisterServer
	id          int
	servers     map[int32]proto.RegisterServer
	timestamp   int32
	portNumber  int
	serverPort1 int
	serverPort2 int
	serverPort3 int
	amount      int
}

var (
	clientPort        = flag.Int("cPort", 0, "client port number")
	serverPort1       = flag.Int("sPort1", 0, "server port number (should match port used for the server)")
	serverPort2       = flag.Int("sPort2", 0, "server port number (should match port used for the server)")
	serverPort3       = flag.Int("sPort3", 0, "server port number (should match port used for the server)")
	clientId          = flag.Int("cId", 0, "client id number")
	serverStream, err proto.RegisterServer
)

//go run server/server.go
//go run client/client.go -cPort 8080 -sPort1 5001 -sPort2 5002 -sPort3 5003 -cId 55

func main() {
	flag.Parse() // Parse the flags to get the port for the client

	cl := &client{
		id:          *clientId,
		servers:     make(map[int32]proto.RegisterServer),
		timestamp:   0,
		portNumber:  *clientPort,
		serverPort1: *serverPort1,
		serverPort2: *serverPort2,
		serverPort3: *serverPort3,
		amount:      50,
	}

	// Wait for input in the client terminal
	scanner := bufio.NewScanner(os.Stdin)

	go registerToServer(cl, *serverPort1, *scanner)
	go registerToServer(cl, *serverPort2, *scanner)
	go registerToServer(cl, *serverPort3, *scanner)
	for {
	}
}

func registerToServer(client *client, serverPort int, scanner bufio.Scanner) {
	//Connect to a server
	serverConnection, _ := connectToServer(serverPort) // This is grpc logic that connects the client to a server

	serverStream, err := serverConnection.JoinServer(context.Background(), &proto.Request{
		Id:   int32(client.id),
		Port: int32(client.portNumber),
	})

	if err != nil {
		log.Printf("Error %s \n", err.Error())
	} else {
		m, e := serverStream.Recv()
		if e != nil {
			log.Printf("Error %s \n", err.Error())
			return
		}
		log.SetFlags(0)
		log.Printf("%d Connected to Server", m.Id)
	}

	go listenOnServer(serverStream, client)

	for scanner.Scan() {
		input := scanner.Text()
		client.timestamp += 1
		//log.Printf("my message: %s", input)

		if input == "bid" {
			serverConnection.Auction(context.Background(), &proto.Bid{
				Amount:  50,
				Comment: "Succes",
			})
		}

	}
}

// This is grpc logic that connects the client to the server
func connectToServer(serverPort int) (proto.RegisterClient, error) {
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to port %v", serverPort)
	} else {
		log.Printf("Connected at port %d\n", serverPort)
	}
	return proto.NewRegisterClient(conn), nil
}

func listenOnServer(serverStream proto.Register_JoinServerClient, client *client) {
	for {
		resp, err := serverStream.Recv()

		if err != nil {
			log.Printf("Error %s", err)
		}
		if resp.Msg == "" {
			log.Printf("%d Connected to Server", resp.Id)

		} else {
			log.Printf("Message from %d: %s", resp.Id, resp.Msg)
		}
	}
}
