package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"strconv"
	"sync"

	proto "github.com/Mastias123/Assignment5.git/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	proto.UnimplementedRegisterServer
	id                            int
	timestamp                     int32
	portNumber                    int
	serverPort1                   int
	serverPort2                   int
	serverPort3                   int
	amount                        int
	myPerseptionOfTheActonsMaxBid int32
}

var servers []proto.RegisterClient
var wg sync.WaitGroup

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
		id: *clientId,
		//servers:     make(map[int32]proto.RegisterServer),
		timestamp:                     0,
		portNumber:                    *clientPort,
		serverPort1:                   *serverPort1,
		serverPort2:                   *serverPort2,
		serverPort3:                   *serverPort3,
		amount:                        50,
		myPerseptionOfTheActonsMaxBid: 0,
	}

	scanner := bufio.NewScanner(os.Stdin)

	go registerToServer(cl, *serverPort1, *scanner)
	go registerToServer(cl, *serverPort2, *scanner)
	go registerToServer(cl, *serverPort3, *scanner)

	go listenOnConsole(cl, *scanner)

	for {
	}
}

func listenOnConsole(client *client, scanner bufio.Scanner) {
	for scanner.Scan() {
		input := scanner.Text()
		client.timestamp += 1

		errorCounter := 0
		bidResultCounter := 0
		maxBid := 0

		if input == "bid" {
			for i := 0; i < len(servers); i++ {
				bidResult, err := servers[i].PlaceBid(context.Background(), &proto.Bid{
					Amount:                        50,
					ClientId:                      int32(client.id),
					ClientPort:                    int32(client.portNumber),
					MyPerseptionOfTheActonsMaxBid: client.myPerseptionOfTheActonsMaxBid,
				})
				if err != nil {
					errorCounter++
				} else {
					bidResultCounter++
					maxBid = int(bidResult.MaxBid)
					client.myPerseptionOfTheActonsMaxBid = int32(maxBid)
				}
			}
			if errorCounter >= bidResultCounter { //Somthing went wrong
				log.Printf("Bid was NOT! accepted") //ToDo denne fejl håndtere både inconsistens mellem serverne og at men ikke har den rigtige perseption af max bid
			} else {
				log.Printf("Succesful bid, max bid is now %d", maxBid) //ToDo evt lav en metode der evt vælger den værdi der forekommer hyppigst
			}

		} else if input == "result" {
			for i := 0; i < len(servers); i++ {
				resultStatus, err := servers[i].Result(context.Background(), &proto.ResultRequest{
					ClientId:   int32(*clientId),
					ClientPort: int32(*clientPort),
				})
				if err != nil { //ToDo Is not used yet

				}

				if client.myPerseptionOfTheActonsMaxBid < resultStatus.MaxBid { // Gets an answer from all the servers and takes the biggest one.
					client.myPerseptionOfTheActonsMaxBid = resultStatus.MaxBid
				}
			}
			log.Printf("Current max Bid: %d", client.myPerseptionOfTheActonsMaxBid)
		}
	}

}

func registerToServer(client *client, serverPort int, scanner bufio.Scanner) {
	//Connect to a server
	serverConnection, _ := connectToServer(serverPort) // This is grpc logic that connects the client to a server

	servers = append(servers, serverConnection)

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

	wg.Add(1)
	wg.Wait()

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

// func listenOnServer(serverStream proto.Register_JoinServerClient, client *client) {
// 	for {
// 		resp, err := serverStream.Recv()

// 		if err != nil {
// 			log.Printf("Error %s", err)
// 		}
// 		if resp.Msg == "" {
// 			log.Printf("%d Connected to Server", resp.Id)

// 		} else {
// 			log.Printf("Message from %d: %s", resp.Id, resp.Msg)
// 		}
// 	}
// }
