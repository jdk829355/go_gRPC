package main

import (
	"flag"
	pb "github.com/jdk829355/go_fssn/server_streaming/serverStreaming"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
)

var (
	addr = flag.String("port", ":50051", "The server port")
)

type server struct {
	pb.UnimplementedServerStreamingServer
}

func makeMessage(s string) pb.Message {
	return pb.Message{Message: s}
}

func (s *server) GetServerResponse(req *pb.Number, stream pb.ServerStreaming_GetServerResponseServer) error {
	messages := []pb.Message{
		makeMessage("message #1"),
		makeMessage("message #2"),
		makeMessage("message #3"),
		makeMessage("message #4"),
		makeMessage("message #5"),
	}
	log.Printf("Server processing gRPC server-streaming {%d}.", req.Value)
	for i := 0; i < 5; i++ {
		if err := stream.Send(&messages[i]); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	s := &server{}
	grpcServer := grpc.NewServer()
	pb.RegisterServerStreamingServer(grpcServer, s)
	log.Printf("Starting server. Listening on port %s.", strings.Split(*addr, ":")[1])
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("%v", err)
	}
}
