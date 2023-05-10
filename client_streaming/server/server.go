package main

import (
	pb "github.com/jdk829355/go_gRPC/client_streaming/ClientStreaming"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

//서버는 무엇을 하는가: client로부터 stream을 받아서 그 stream의 개수를 반환

type server struct {
	pb.UnimplementedClientStreamingServer
}

// GetServerResponse 클라이언트로부터 받은 스트림에 대한 응답
func (s *server) GetServerResponse(stream pb.ClientStreaming_GetServerResponseServer) error {
	counter := 0
	log.Println("Server processing gRPC client-streaming.")
	for {
		// stream의 한 요소를 Recv() 통하여 받음
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error when reading client request stream: %v", err)
		}
		// 개수를 카운트하기 위해 반복
		counter++
	}
	// 스트림으로 응답을 주고 스트림 종료
	err := stream.SendAndClose(&pb.Number{Value: int32(counter)})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("%v", err)
	}
	s := grpc.NewServer()
	pb.RegisterClientStreamingServer(s, &server{})
	log.Println("Starting server. Listening on port 50051.")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
