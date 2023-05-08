package main

import (
	"flag"
	pb "github.com/jdk829355/go_fssn/server_streaming/serverStreaming"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
)

// 명령행 인자로 소켓 주소를 받음
var (
	addr = flag.String("port", ":50051", "The server port")
)

// unary와 마찬가지로 구현할 함수가 있는 서버 구조체를 임베드
type server struct {
	pb.UnimplementedServerStreamingServer
}

// 문자열을 value로 가지는 Message 인스턴스 반환
func makeMessage(s string) pb.Message {
	return pb.Message{Message: s}
}

// GetServerResponse proto file에서 정의했던 함수 구현
func (s *server) GetServerResponse(req *pb.Number, stream pb.ServerStreaming_GetServerResponseServer) error {
	// 순차적으로 보낼 메시지
	messages := []pb.Message{
		makeMessage("message #1"),
		makeMessage("message #2"),
		makeMessage("message #3"),
		makeMessage("message #4"),
		makeMessage("message #5"),
	}
	// request에 있는 Value를 출력
	log.Printf("Server processing gRPC server-streaming {%d}.", req.Value)
	// 에러 체크를 하며 stream.Send를 통해 메시지를 보낸다.
	for i := 0; i < 5; i++ {
		if err := stream.Send(&messages[i]); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// 나머지 과정은 unary와 동일
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
