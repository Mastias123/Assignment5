package main

import (
	"bufio"
	"os"
	"strconv"

	proto "github.com/Mastias123/Assignment5.git/grpc"
)

type client struct {
	proto.UnimplementedRegisterServer
	id        int32
	servers   map[int32]proto.RegisterServer
	timestamp int32
}

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32) //Could also use flags
	ownPort := int32(arg1) + 5000

	cl := &client{
		id:        ownPort,
		servers:   make(map[int32]proto.RegisterServer),
		timestamp: 0,
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "bid" {

			//strconv.ParseInt(scanner.Text()
			go cl.bid()
		}
	}

}

// amount int32)
func (c *client) bid() {

}
