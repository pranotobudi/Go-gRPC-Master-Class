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
	// doServerStreaming(c)
	// doClientStreaming(c)
	doBiDiStreaming(c)

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

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Println("failed")
	}
	var requests = []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 4,
		},
	}
	for _, request := range requests {
		stream.Send(request)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Println("failed receiving data")
	}
	log.Println("average: ", res.GetResult())
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Println("failed connecting the server...")
	}
	//send messages
	var numbers = []int64{6, 1, 2, 3, 4, 5}
	waitc := make(chan struct{})
	// var wg sync.WaitGroup
	go func() {
		for _, val := range numbers {
			time.Sleep(100 * time.Millisecond)
			req := &calculatorpb.FindMaximumRequest{
				Number: val,
			}
			// wg.Add(1)
			stream.Send(req)
			// wg.Done()
		}
		stream.CloseSend()
	}()

	//receive messages
	go func() {
		for {
			// wg.Add(1)
			res, err := stream.Recv()
			// wg.Done()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("failed to receive the message")
			}
			fmt.Println("result: ", res.GetResult())
		}
		close(waitc)
	}()

	// wg.Wait()
	<-waitc
}
