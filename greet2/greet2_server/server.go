package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/pranotobudi/Go-gRPC-Master-Class/greet2/greet2pb"
	"google.golang.org/grpc"
)

type server struct {
	greet2pb.UnimplementedGreetServiceServer
}

func (s *server) Greet(ctx context.Context, req *greet2pb.GreetRequest) (*greet2pb.GreetResponse, error) {
	fmt.Printf("Greet function is called with %v", req)
	firstName := req.GetGreeting().FirstName
	// firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().LastName
	result := "Hello " + firstName + " " + lastName
	response := greet2pb.GreetResponse{
		Result: result,
	}
	return &response, nil
}

func (s *server) GreetManyTimes(req *greet2pb.GreetManyTimeRequest, stream greet2pb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().FirstName
	for i := 0; i < 5; i++ {
		result := "#" + strconv.Itoa(i) + "Hello, " + firstName
		log.Println("server result: ", result)
		res := &greet2pb.GreetManyTimeResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(time.Millisecond * 1000)
	}
	return nil
}

func (s *server) LongGreet(stream greet2pb.GreetService_LongGreetServer) error {
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
		}
		firstName := req.GetGreeting().FirstName
		result += "Hello" + firstName + "!"
	}
	res := &greet2pb.LongGreetResponse{
		Result: result,
	}
	stream.SendAndClose(res)
	return nil
}

func (s *server) GreetEveryone(stream greet2pb.GreetService_GreetEveryoneServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream failed %v", err)
		}
		result := "Hello, " + req.GetGreeting().FirstName
		response := &greet2pb.GreetEveryoneResponse{
			Result: result,
		}
		err = stream.Send(response)
		if err != nil {
			log.Fatalf("error while sending data: %v", err)
		}
	}
	return nil
}

func main() {
	fmt.Println("bismillah")

	s := grpc.NewServer()
	greet2pb.RegisterGreetServiceServer(s, &server{})

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
