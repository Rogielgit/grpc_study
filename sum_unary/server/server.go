package main

import (
	"context"
	"example.com/grpc_study/sum_unary/sum_proto"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type payload struct {
	sum_proto.UnimplementedSumServiceServer
}

func (p *payload) Sum(cxt context.Context, req *sum_proto.SumRequest) (*sum_proto.SumResponse, error) {

	result := req.Values.N1 + req.Values.N2

	return &sum_proto.SumResponse{
		Result: result,
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
	sum_proto.RegisterSumServiceServer(s, &payload{})

	if err = s.Serve(lis); err != nil {
		log.Fatal("Failed on server execution")
	}
}