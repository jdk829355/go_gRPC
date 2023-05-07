package main

import (
	"context"
	pb "github.com/jdk829355/go_fssn/client_streaming/ClientStreaming"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func makeMessage(s string) *pb.Message {
	return &pb.Message{Message: s}
}

func main() {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(":50051", opts)
	if err != nil {
		log.Fatalf("Error connecting: %v \n", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}(conn)
	c := pb.NewClientStreamingClient(conn)
	getServerResponse(c)
}

func getServerResponse(c pb.ClientStreamingClient) {
	req := []*pb.Message{
		makeMessage("message #1"),
		makeMessage("message #2"),
		makeMessage("message #3"),
		makeMessage("message #4"),
		makeMessage("message #5"),
	}
	stream, err := c.GetServerResponse(context.Background())
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, r := range req {
		log.Printf("[client to server] %s", r.GetMessage())
		err := stream.Send(r)
		if err != nil {
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("[server to client] %d", res.GetValue())
}
