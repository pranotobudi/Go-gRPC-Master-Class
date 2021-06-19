package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/pranotobudi/Go-gRPC-Master-Class/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Hello I'm a client")

	// tls := false
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	// doUnary(c)
	doServerClient(c)

}

func doUnary(c greetpb.GreetServiceClient) {
	req := greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "bud",
			LastName:  "sas",
		},
	}
	res, err := c.Greet(context.Background(), &req)
	if err != nil {
		log.Println("cant get response")
	}
	fmt.Printf("Response \n %v", res)

}

func doServerClient(c greetpb.GreetServiceClient) {
	req := greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "bud",
			LastName:  "sas",
		},
	}
	stream, err := c.GreetManyTimes(context.Background(), &req)
	if err != nil {
		log.Println("cant get response")
	}
	for {
		msg, err := stream.Recv()
		fmt.Printf("Response %v \n", msg)
		fmt.Printf("Response %v \n", msg.GetResult())
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("error while reading stream")
		}
	}

}
