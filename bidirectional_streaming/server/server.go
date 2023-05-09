package main

import (
	pb "github.com/jdk829355/go_fssn/bidirectional_streaming/bidirectionalStreaming"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedBidirectionalServer
}

func (s *server) GetServerResponse(stream pb.Bidirectional_GetServerResponseServer) error {
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else {
			log.Println("Server processing gRPC bidirectional streaming.")
			err := stream.Send(message)
			if err != nil {
				return err
			}
		}

	}

}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("%v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBidirectionalServer(s, &server{})
	log.Println("Starting server. Listening on port 50051.")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("%v", err)
	}
}
