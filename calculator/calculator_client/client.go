package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/pranotobudi/Go-gRPC-Master-Class/calculator/calculatorpb"
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
	c := calculatorpb.NewCalculatorServiceClient(cc)
	// doUnary(c)
	doServerStreaming(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := calculatorpb.SumRequest{
		Operand: &calculatorpb.Operand{
			FirstElement:  1,
			SecondElement: 2,
		},
	}

	res, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Println("cant get response")
	}
	log.Printf("Response Sum \n %v", res)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC Client...")
	req := calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), &req)
	if err != nil {
		log.Println("cant get response")
	}
	for {
		time.Sleep(1000 * time.Millisecond)
		msg, err := stream.Recv()
		fmt.Println("msg", msg, "err:", err)

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("error processing data")
		}
		log.Println("response: ", msg.GetResult())

	}

}
