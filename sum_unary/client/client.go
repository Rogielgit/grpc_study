package main

import (
	"context"
	"example.com/grpc_study/sum_unary/sum_proto"
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

	c := sum_proto.NewSumServiceClient(conn)
	doUnaryConnection(c)

}

func doUnaryConnection(c sum_proto.SumServiceClient) {
	sum := sum_proto.Sum{N1: 5,N2: 10}
	req := sum_proto.SumRequest{
		Values: &sum,
	}

	resp, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Fatal("error while calling Greet RPC: %v", err)
	}
	log.Printf("response from server: %v", resp.Result)
}
