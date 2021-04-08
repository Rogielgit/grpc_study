package main

import (
	"example.com/grpc_study/greet/greetpb"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type  server struct {
	greetpb.UnimplementedGreetServiceServer
}

func main()  {

	fmt.Print("Starting the application!!\n")

	lis, err := net.Listen("tcp", "0.0.0.0:5154")
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	_ = s
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatal("Failed on server execution")
	}

}








