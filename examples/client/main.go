package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/ulule/ipfix/proto"
	"google.golang.org/grpc"
)

func main() {
	var serverAddr string
	var ipAddress string

	flag.StringVar(&serverAddr, "server.addr", "127.0.0.1:33001", "rpc server addr")
	flag.StringVar(&ipAddress, "ip", "127.0.0.1", "ip address")
	flag.Parse()

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewIpfixClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	loc, err := c.GetLocation(ctx, &proto.GetLocationRequest{IpAddress: ipAddress})
	if err != nil {
		log.Fatalf("could not retrieve result: %v", err)
	}

	log.Printf("Retrieved location: %+v", loc)
}
