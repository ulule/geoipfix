package main

import (
	"context"
	"flag"
	"log"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/ulule/ipfix/proto"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func main() {
	var serverAddr string
	var ipAddress string
	var maxRetry uint

	flag.StringVar(&serverAddr, "server-addr", "127.0.0.1:33001", "rpc server addr")
	flag.StringVar(&ipAddress, "ip", "127.0.0.1", "ip address")
	flag.UintVar(&maxRetry, "retries", 3, "max retries")
	flag.Parse()

	tr := trace.New("ipfix.Client", "GetLocation")
	defer tr.Finish()

	conn, err := grpc.Dial(serverAddr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewIpfixClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ctx = trace.NewContext(ctx, tr)
	loc, err := c.GetLocation(ctx, &proto.GetLocationRequest{IpAddress: ipAddress},
		grpc_retry.WithMax(maxRetry),
		grpc_retry.WithPerRetryTimeout(1*time.Second),
		grpc_retry.WithCodes(codes.Unavailable))

	if err != nil {
		log.Fatalf("could not retrieve result: %v", err)
	}

	log.Printf("Retrieved location: %+v", loc)
}
