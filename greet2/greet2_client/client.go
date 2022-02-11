package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/pranotobudi/Go-gRPC-Master-Class/greet2/greet2pb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("I'm client..")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect:%v", err)
	}
	defer conn.Close()
	client := greet2pb.NewGreetServiceClient(conn)
	fmt.Println("created client: %v", client)
	// doUnary(client)
	// doServerStreaming(client)
	// doClientStreaming(client)
	doBidirectionalStreaming(client)
}

func doUnary(client greet2pb.GreetServiceClient) {
	log.Println("doUnary..")
	req := &greet2pb.GreetRequest{
		Greeting: &greet2pb.Greeting{
			FirstName: "pranoto",
			LastName:  "budi",
		},
	}
	res, err := client.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("response failed")
	}
	log.Println("response: ", res)

}

func doServerStreaming(client greet2pb.GreetServiceClient) {
	log.Println("doServerStreaming..")
	req := &greet2pb.GreetManyTimeRequest{
		Greeting: &greet2pb.Greeting{
			FirstName: "pranoto",
			LastName:  "budi",
		},
	}
	resStream, err := client.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("response failed")
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}
		log.Println("Message: ", msg.GetResult())
	}
}

func doClientStreaming(client greet2pb.GreetServiceClient) {
	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("failed to get stream %v", err)
	}
	for i := 0; i < 5; i++ {
		req := &greet2pb.LongGreetRequest{
			Greeting: &greet2pb.Greeting{
				FirstName: "pranoto",
				LastName:  "budi",
			},
		}
		stream.Send(req)
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to get response %v", err)
	}
	log.Println("response: ", response.String())
}

func doBidirectionalStreaming(client greet2pb.GreetServiceClient) {
	stream, err := client.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("failed to get client stream: %v", err)
	}
	waitch := make(chan struct{})
	go func() {
		requests := []*greet2pb.GreetEveryoneRequest{
			{
				Greeting: &greet2pb.Greeting{
					FirstName: "budi1",
				},
			},
			{
				Greeting: &greet2pb.Greeting{
					FirstName: "budi2",
				},
			},
			{
				Greeting: &greet2pb.Greeting{
					FirstName: "budi3",
				},
			},
			{
				Greeting: &greet2pb.Greeting{
					FirstName: "budi4",
				},
			},
		}
		for _, req := range requests {
			stream.Send(req)
			time.Sleep(time.Microsecond * 1000)
		}
		err = stream.CloseSend()
		if err != nil {
			log.Fatalf("failed to close client stream %v", err)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitch)
				break
			}
			if err != nil {
				log.Fatalf("failed to receive response: %v", err)
			}
			log.Println("result: ", res.GetResult())
		}
	}()

	<-waitch
}
