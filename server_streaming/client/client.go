package main

import (
	"context"
	"flag"
	pb "github.com/jdk829355/go_fssn/server_streaming/serverStreaming"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

var (
	serverAddr = flag.String("server_addr", "localhost:50051", "The server address with port")
	value      = flag.Int("value", 5, "value to calculate")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}(conn)

	client := pb.NewServerStreamingClient(conn)
	stream, err := client.GetServerResponse(context.Background(), &pb.Number{Value: int32(*value)})
	if err != nil {
		log.Fatalf("GetServerResponse - %v", err)
	}

	for {
		content, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("GetServerResponse stream - %v", err)
		}

		log.Printf("[server to client] %s", content.GetMessage())
	}
}
