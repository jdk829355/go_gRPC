package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/jdk829355/go_fssn/unary/unaryService"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultvalue = 2
)

var (
	addr  = flag.String("addr", "localhost:50051", "the address to connect to")
	value = flag.Int("value", defaultvalue, "value to calculate")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUnaryServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.MyFunction(ctx, &pb.MyNumber{Value: int32(*value)})
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	log.Printf("gRPC result:: %d", int(r.GetValue()))
}
