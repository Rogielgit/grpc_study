package main

import (
	"context"
	"example.com/grpc_study/greet/protobf"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct {
	protobf.UnimplementedGreetServiceServer
}

func (s *server) GreetManyTimes(req *protobf.GreetManyTimesRequest, stream protobf.GreetService_GreetManyTimesServer) error {
	fmt.Printf("Greet server streaming func was called!!")
	firstName := req.Greeting.GetFirstName()
	for i := 0; i < 20; i++ {
		stringResp := "Hello " + firstName + " number " + strconv.Itoa(i)

		res := protobf.GreetManyTimesResponse{
			Result: stringResp,
		}

		stream.Send(&res)
		time.Sleep(1000 * time.Microsecond)
	}

	return nil
}

func (s *server) Greet(ctx context.Context, req *protobf.GreetRequest) (*protobf.GreetResponse, error) {
	fmt.Printf("Greet func was called!!")

	fmt.Printf(req.Greeting.GetFirstName())
	fmt.Println(req.Greeting.GetLastName())
	res := "Hello " + req.Greeting.GetFirstName()

	return &protobf.GreetResponse{
		Result: res,
	}, nil
}

func main() {
	fmt.Print("Starting the application!!\n")

	lis, err := net.Listen("tcp", "0.0.0.0:5051")
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	_ = s
	protobf.RegisterGreetServiceServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatal("Failed on server execution")
	}
}
