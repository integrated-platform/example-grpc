// main.go
package main

import (
	"context"
	"log"
	"net"
	"os"
	"sync"
	"time"

	pb "example-grpc/proto/out"

	"google.golang.org/grpc"
)

// 서버 구현
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello 구현
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func startServer(wg *sync.WaitGroup) {
	defer wg.Done()

	const port = ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// 클라이언트 구현
func startClient() {
	const address = "localhost:50051"
	defaultName := "world"

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go startServer(&wg)

	// 클라이언트를 잠시 후 실행
	time.Sleep(1 * time.Second)
	startClient()

	wg.Wait()
}
