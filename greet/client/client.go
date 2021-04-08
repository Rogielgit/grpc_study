package main

import (
	"context"
	"example.com/grpc_study/greet/protobf"
	"fmt"
	"google.golang.org/grpc"
	"io"
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
	doServerStream(c)
	//doUnaryConnection(c)

}

func doServerStream(c protobf.GreetServiceClient) {
	g := protobf.Greeting{FirstName: "Test_First_name", LastName: "Test_last_name"}
	req := protobf.GreetManyTimesRequest{
		Greeting: &g,
	}

	resp, err := c.GreetManyTimes(context.Background(), &req)
	if err != nil {
		log.Fatal("error while calling Greet streammig RPC: %v", err)

	}



	for {
		streamRes, err := resp.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("error while reading server stream RPC: %v", err)
		}

		log.Printf("Response from server: %v", streamRes.GetResult())
	}
}

func doUnaryConnection(c protobf.GreetServiceClient) {
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
