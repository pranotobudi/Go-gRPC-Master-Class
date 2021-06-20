package main

import (
	"context"
	"fmt"
	"io"
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

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, res_stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
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
			res_stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
			number = number / k // divide N by k so that we have the rest of the number left.
		} else {
			k = k + 1
		}

	}
	return nil
}

//Client Streaming
func (s *server) ComputeAverage(req_stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	var total int64
	var counter int64
	var avg int64
	for {
		msg, err := req_stream.Recv()
		if err == io.EOF {
			avg = total / counter
			break
		}
		if err != nil {
			log.Println("failed processing message")
			return err
		}
		total += msg.GetNumber()
		counter++
		log.Println("number: ", msg.GetNumber())
	}
	res := &calculatorpb.ComputeAverageResponse{
		Result: avg,
	}
	return req_stream.SendAndClose(res)
}

//BiDi Streaming
func (s *server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	var numbers []int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		number := req.GetNumber()
		fmt.Println("number: ", number)
		numbers = append(numbers, number)
		max := findMax(numbers)
		res := &calculatorpb.FindMaximumResponse{
			Result: max,
		}
		stream.Send(res)
	}
	return nil
}

func findMax(numbers []int64) int64 {
	var max int64
	for _, val := range numbers {
		if val > max {
			max = val
		}
	}
	return max
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
