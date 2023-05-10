package main

import (
	"context"
	"flag"
	pb "github.com/jdk829355/go_gRPC/server_streaming/serverStreaming"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

// 명령행 인자로 소켓 주소와 함수에 넣을 인자를 받는다
var (
	serverAddr = flag.String("server_addr", "localhost:50051", "The server address with port")
	value      = flag.Int("value", 5, "value to calculate")
)

func main() {
	flag.Parse()
	// grpc서버에 연결
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	// 연결 해제 예약
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}(conn)

	// 클라이언트 생성
	client := pb.NewServerStreamingClient(conn)
	// 함수 호출
	stream, err := client.GetServerResponse(context.Background(), &pb.Number{Value: int32(*value)})
	if err != nil {
		log.Fatalf("GetServerResponse - %v", err)
	}
	// stream을 받기
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
