package main

import (
	"context"
	"example.com/grpc_study/unary/protobf"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func main() {

	fmt.Print("Hello I'm a client!!\n")

	conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("could not connect to server %v", err)
	}

	defer conn.Close()

	c := protobf.NewGreetServiceClient(conn)
	g := protobf.Greeting{FirstName: "Test_First_name", LastName: "Test_last_name"}
	req := protobf.GreetRequest{
		Greeting: &g,
	}

	resp, err := c.Greet(context.Background(), &req)
	if err != nil {
		log.Fatal("error while calling Greet RPC: %v", err)
	}
	log.Printf("response from greet: %v", resp.Result)
}
