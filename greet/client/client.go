package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"example.com/grpc_study/greet/protobf"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	pemCA, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemCA) {
		return nil, fmt.Errorf("Failed to add server CA's certificate!!")
	}

	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	fmt.Print("Hello I'm a client!!\n")

	tlsCred, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: %v", err)
	}

	conn, err := grpc.Dial("0.0.0.0:5051", grpc.WithTransportCredentials(tlsCred))
	if err != nil {
		log.Fatal("could not connect to server %v", err)
	}
	defer conn.Close()

	c := protobf.NewGreetServiceClient(conn)
	doClientStreaming(c)
	//doServerStream(c)
	//doUnaryConnection(c)
}

func doClientStreaming(c protobf.GreetServiceClient) {
	fmt.Println("Starting client streaming")

	stream, err := c.ClientStreamGreet(context.Background())
	if err != nil {
		log.Fatal("error while calling client streammig RPC: %v", err)
	}

	for i := 0; i < 10; i++ {
		g := protobf.Greeting{FirstName: "Test_First_name" + strconv.Itoa(i), LastName: "Test_last_name" + strconv.Itoa(i)}
		req := protobf.ClientStreamGreetRequest{
			Greeting: &g,
		}

		fmt.Printf("seding the request %v\n", req)
		stream.Send(&req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("error while calling client streammig RPC: %v", err)
	}
	fmt.Printf("response:  %v", res)

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
