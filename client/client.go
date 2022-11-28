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
	id                      int
	timestamp               int32
	portNumber              int
	serverPort1             int
	serverPort2             int
	serverPort3             int
	amount                  int
	myPerseptionOfTheMaxBid int32
	hasMaxBidId             int32
}

var servers []proto.RegisterClient
var wg sync.WaitGroup
var wt sync.WaitGroup

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
		timestamp:               0,
		portNumber:              *clientPort,
		serverPort1:             *serverPort1,
		serverPort2:             *serverPort2,
		serverPort3:             *serverPort3,
		amount:                  50,
		myPerseptionOfTheMaxBid: 0,
	}

	scanner := bufio.NewScanner(os.Stdin)

	go registerToServer(cl, *serverPort1, *scanner)
	go registerToServer(cl, *serverPort2, *scanner)
	go registerToServer(cl, *serverPort3, *scanner)

	go listenOnConsole(cl, *scanner)

	wt.Add(1)
	wt.Wait()
}

func listenOnConsole(client *client, scanner bufio.Scanner) {
	for scanner.Scan() {
		input := scanner.Text()
		client.timestamp += 1

		if input == "bid" {

			maxBid := 0
			errorCounter := 0
			bidResultCounter := 0

			for i := 0; i < len(servers); i++ {
				//front enden skal time out'e serveren
				bidResult, err := servers[i].PlaceBid(context.Background(), &proto.Bid{
					Amount:                        50,
					ClientId:                      int32(client.id),
					ClientPort:                    int32(client.portNumber),
					MyPerseptionOfTheActonsMaxBid: client.myPerseptionOfTheMaxBid,
				})
				if err != nil {
					errorCounter++
				} else {
					bidResultCounter++
					maxBid = int(bidResult.MaxBid)
					client.myPerseptionOfTheMaxBid = int32(maxBid)
					client.hasMaxBidId = bidResult.MaxBidId //Burde rykkes
				}
			}
			if errorCounter >= bidResultCounter { //Somthing went wrong
				log.Printf("Bid was NOT! accepted. Please check the result") //ToDo denne fejl håndtere både inconsistens mellem serverne og at man ikke har den rigtige perseption af max bid
			} else {
				log.Printf("Succesful bid, max bid is now %d", maxBid) //ToDo evt lav en metode der evt vælger den værdi der forekommer hyppigst
			}

			if maxBid == 200 {
				log.Printf("Auction is over")
			}

		} else if input == "result" {

			for i := 0; i < len(servers); i++ {
				resultStatus, err := servers[i].Result(context.Background(), &proto.ResultRequest{
					ClientId:   int32(client.id),
					ClientPort: int32(client.portNumber),
				})

				if err != nil { //ToDo Is not used yet, is there to fix compile error

				}

				if client.myPerseptionOfTheMaxBid < resultStatus.MaxBid { // Gets an answer from all the servers and takes the biggest one.
					client.myPerseptionOfTheMaxBid = resultStatus.MaxBid
					client.hasMaxBidId = resultStatus.Id
				}
			}
			//log.Printf("clientId %d clientHasMaxBidId %d", client.id, client.hasMaxBidId)
			if int32(client.id) == client.hasMaxBidId {
				log.Printf("You have the current max Bid: %d", client.myPerseptionOfTheMaxBid)
			} else {
				log.Printf("Id: %d, has the current max Bid: %d", client.hasMaxBidId, client.myPerseptionOfTheMaxBid)
			}

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
