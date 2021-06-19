package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/pranotobudi/Go-gRPC-Master-Class/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	firstElmt := req.GetOperand().GetFirstElement()
	secondElmt := req.GetOperand().GetSecondElement()
	result := firstElmt + secondElmt
	res := &calculatorpb.SumResponse{
		Result: result,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function was invoked with %v\n", req)

	number := req.GetNumber()
	k := int64(2)
	fmt.Println("1: ", number)
	for number > 1 {
		// fmt.Println("number: ", number)
		if number%k == 0 { // if k evenly divides into N
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: k, // this is a factor
			}
			stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
			number = number / k // divide N by k so that we have the rest of the number left.
		} else {
			k = k + 1
		}

	}
	return nil
}

func main() {
	fmt.Println("server running...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
