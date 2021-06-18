package main

import (
	"context"
	"fmt"
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
	fmt.Printf("Response %v", res)

}
