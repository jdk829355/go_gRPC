package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/jdk829355/go_fssn/unary/unaryService"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultvalue = 2
)

var (
	// 명령행 인자로 연결할 소켓 주소와 계산할 값을 입력받음
	addr  = flag.String("addr", "localhost:50051", "the address to connect to")
	value = flag.Int("value", defaultvalue, "value to calculate")
)

func main() {
	flag.Parse()

	// 서버에 연결하는 gRPC 클라이언트의 연결을 설정
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// 함수가 끝나기 전 연결 해제
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}(conn)

	// 클라이언트에서 gRPC stub을 생성합니다(MyFunction을 멤버 함수로 갖고 있음)
	c := pb.NewUnaryServiceClient(conn)

	// context 객체로 1초의 timeout 설정
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	// cancel(): context와 관련된 자원을 종료
	// Canceling this context releases resources associated with it,
	// so code should call cancel as soon as the operations running in this Context complete
	defer cancel()
	r, err := c.MyFunction(ctx, &pb.MyNumber{Value: int32(*value)})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("gRPC result:: %d", int(r.GetValue()))
}
