package main

import (
	"context"
	pb "github.com/jdk829355/go_gRPC/bidirectional_streaming/bidirectionalStreaming"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

func makeMessage(s string) *pb.Message {
	return &pb.Message{Message: s}
}

func main() {
	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(":50051", opt)
	if err != nil {
		log.Fatalf("%v", err)
	}

	client := pb.NewBidirectionalClient(conn)
	stream, err := client.GetServerResponse(context.Background())
	if err != nil {
		log.Fatalf("%v", err)
	}

	waitc := make(chan bool)
	//ctx := stream.Context()

	messages := []*pb.Message{
		makeMessage("message #1"),
		makeMessage("message #2"),
		makeMessage("message #3"),
		makeMessage("message #4"),
		makeMessage("message #5"),
	}

	go func() {
		for {
			message, err := stream.Recv()

			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				log.Printf("line 50: %v", err)
			}

			log.Printf("[server to client] %s", message.GetMessage())
		}
	}()

	for _, message := range messages {
		log.Printf("[client to server] %s", message.GetMessage())
		if err := stream.Send(message); err != nil {
			log.Fatalf("%v", err)
		}
	}
	if err := stream.CloseSend(); err != nil {
		log.Println(err)
	}

	<-waitc
	err = conn.Close()
	if err != nil {
		return
	}
}
