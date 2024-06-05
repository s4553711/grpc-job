package simple_server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	pb "github.com/s4553711/grpc-jobs/hello"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedHelloServiceServer
	ch chan int
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloResponse{Reply: "Hello " + in.GetName()}, nil
}

// RegNode is the place for slave node to register and add to available host list
func (s *server) RegNode(ctx context.Context, in *pb.RegReq) (*pb.RegRep, error) {
	log.Printf("Received slave %s", in.GetHostname())
	return &pb.RegRep{Reply: "Hello "+in.GetHostname()}, nil
}

func (s *server) ReqJob(ctx context.Context, in *pb.JobReq) (*pb.HelloResponse, error) {
	log.Printf("Ready to send job: %v", in.GetCommand())
	cl := CreateNewAGClient("localhost:50052")
	cl.ExecComm(in.GetCommand())
	return &pb.HelloResponse{Reply: "exec " + in.GetCommand()}, nil
}

func (s *server) Terminate(ctx context.Context, in *pb.Empty) (*pb.HelloResponse, error) {
	log.Printf("Ready to terminate server")
	s.ch <- 1
	return &pb.HelloResponse{Reply: "done"}, nil
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

	ch := make(chan int, 0)
	gsrv := server{ ch: ch }
	sg := grpc.NewServer()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		s := <-ch
		log.Printf("got signal %v, attempting graceful shutdown", s)
		//cancel()
		sg.GracefulStop()
		wg.Done()
	}()

	//pb.RegisterHelloServiceServer(sg, &server{})
	pb.RegisterHelloServiceServer(sg, &gsrv)
	log.Printf("server listening at %v", lis.Addr())
	if err := sg.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	wg.Wait()
	log.Println("Server shutdown")
}
