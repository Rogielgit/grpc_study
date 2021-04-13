package main

import (
	"context"
	"example.com/grpc_study/poc/poc_proto"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: true,            // Allow pings even when there are no active streams
}

var kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
}

// server implements POC.
type server struct {
	poc_proto.UnimplementedCheckServiceServer
}

func (s *server) HasWork(context.Context, *poc_proto.CheckWorkRequest) (*poc_proto.CheckWorkResponse, error) {
	fmt.Print("checking if server has work to client..\n")
	var hasWork poc_proto.CheckWork
	if hasWorkToClient() {
		hasWork = poc_proto.CheckWork{
			HasWork: true,
		}
		return &poc_proto.CheckWorkResponse{
			Check: &hasWork,
		}, nil
	}

	return &poc_proto.CheckWorkResponse{
		Check: &poc_proto.CheckWork{HasWork: false},
	}, nil
}

func (s *server) DoWork(stream poc_proto.CheckService_DoWorkServer) error {
	fmt.Println("server has work to client..")
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&poc_proto.DoWorkResponse{
				Message: "received all values",
			})
		}
		if err != nil {
			return err
		}

		fmt.Printf("value received: %v\n", res.Network)
	}

	return nil
}

func hasWorkToClient() bool {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(1000) > 800 {
		return true
	}
	return false
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:5051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
	poc_proto.RegisterCheckServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
