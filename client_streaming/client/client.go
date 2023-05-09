package main

import (
	"context"
	pb "github.com/jdk829355/go_fssn/client_streaming/ClientStreaming"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// makeMessage: string을 받아 그걸 메시지로 가지는 Message 인스턴스의 주소값 반환
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

// getServerResponse: protoc에 의해 자동 생성된 코드에 있는 함수를 직접 호출하는 것이 아니라 stream을 보내는 함수를 따로 만듦
// main 함수에서 만들어진 클라이언트 인스턴스를 인자로 한 함수
func getServerResponse(c pb.ClientStreamingClient) {
	req := []*pb.Message{
		makeMessage("message #1"),
		makeMessage("message #2"),
		makeMessage("message #3"),
		makeMessage("message #4"),
		makeMessage("message #5"),
	}
	// 인자로 받은 클라이언트에 있는 protoc의 함수를 호출하면 stream을 반환함
	stream, err := c.GetServerResponse(context.Background())
	if err != nil {
		log.Fatalf("%v", err)
	}

	// stream.Send()로 메시지 보내기
	for _, r := range req {
		log.Printf("[client to server] %s", r.GetMessage())
		err := stream.Send(r)
		if err != nil {
			return
		}
	}

	// 다 보내면 응답을 받기
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("[server to client] %d", res.GetValue())
}
