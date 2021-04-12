package main

import (
	"context"
	"crypto/tls"
	"example.com/grpc_study/greet/protobf"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct {
	protobf.UnimplementedGreetServiceServer
}

func (s *server) ClientStreamGreet(stream protobf.GreetService_ClientStreamGreetServer) error {
	fmt.Println("ClientStreamGreet was invoked")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&protobf.ClientStreamGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			return nil
		}

		firstName := req.Greeting.GetFirstName()
		result += "Hello " + firstName + "!"
	}

	return nil
}

func (s *server) GreetManyTimes(req *protobf.GreetManyTimesRequest, stream protobf.GreetService_GreetManyTimesServer) error {
	fmt.Printf("Greet server streaming func was called!!\n")
	firstName := req.Greeting.GetFirstName()

	for i := 0; i < 20; i++ {
		stringResp := "Hello " + firstName + " number " + strconv.Itoa(i)

		res := protobf.GreetManyTimesResponse{
			Result: stringResp,
		}

		stream.Send(&res)
		time.Sleep(2000 * time.Millisecond)
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

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	fmt.Print("Starting the application!!\n")

	lis, err := net.Listen("tcp", "0.0.0.0:5051")
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}

	tlsCred, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials!")
	}

	s := grpc.NewServer(grpc.Creds(tlsCred))
	protobf.RegisterGreetServiceServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatal("Failed on server execution")
	}
}
