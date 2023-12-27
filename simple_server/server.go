package simple_server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "github.com/s4553711/grpc-jobs/hello"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedHelloServiceServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloResponse{Reply: "Hello " + in.GetName()}, nil
}

type Gserver struct {
	port int
}

func CreateNew(port int) *Gserver {
	return &Gserver{port: port}
}

func (s *Gserver) RunServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	sg := grpc.NewServer()
	pb.RegisterHelloServiceServer(sg, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := sg.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
