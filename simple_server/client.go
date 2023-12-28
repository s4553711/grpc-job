package simple_server

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "github.com/s4553711/grpc-jobs/hello"
)

type Gclient struct {
	addr string
	conn *grpc.ClientConn
	sc pb.HelloServiceClient
}

func CreateNewClient(addr string) *Gclient {
	return &Gclient{addr: addr}
}

func (c *Gclient) Connect() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c.conn = conn
	c.sc = pb.NewHelloServiceClient(conn)
}

func (c *Gclient) RegNode(host string, port int) {
	c.Connect()
	defer c.conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.sc.RegNode(ctx, &pb.RegReq{Hostname: host, Port: uint32(port)})
	if err != nil {
		log.Fatal("could not register node: %v", err)
	}
	log.Printf("Ack from RegNode: %s", r.GetReply())
}

func (c *Gclient) SayHello() {
	c.Connect()
	defer c.conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.sc.SayHello(ctx, &pb.HelloRequest{Name: "Ivy"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Ack from Greeting: %s", r.GetReply())
}
