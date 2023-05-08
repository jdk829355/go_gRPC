package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/jdk829355/go_fssn/unary/unaryService"
	"google.golang.org/grpc"
)

// --port에서 포트넘버를 명령행 인자로 받음 (기본값: 50051)
var (
	port = flag.Int("port", 50051, "The server port")
)

// protoc에 정의된 함수를 멤버 함수로 갖는 구조체를 임베드
type server struct {
	pb.UnimplementedUnaryServiceServer
}

// MyFunction proto 파일 내 정의한 rpc 함수 이름에 대응하는 함수 작성
func (s *server) MyFunction(ctx context.Context, in *pb.MyNumber) (*pb.MyNumber, error) {
	// 사전에 정의한 MyFunction_unary(정수의 제곱 함수)를 MyFunction에 구현
	result := pb.MyFunction_unary(int(in.GetValue()))
	return &pb.MyNumber{Value: int32(result)}, nil
}

func main() {
	flag.Parse()

	// tcp 바인딩
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 요청을 받지 않고 아무런 서비스도 등록되지 않은 grpc서버 구조체(Server{}) 반환
	s := grpc.NewServer()

	// 서비스가 없는 s라는 grpc서버에 MyFunction()이 구현된 server{}를 등록
	pb.RegisterUnaryServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	// net.Listener인 lis로 들어오는 연결 요청을 받아들이고 각 연결마다 goroutine으로 비동기 처리
	// error를 반환하며 이 값이 nil이 아니면 os.Exit(1)로 프로그램 종료
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
