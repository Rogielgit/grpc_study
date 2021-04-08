package main

import (
	"example.com/grpc_study/greet/greetpb"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func main()  {

	fmt.Print("Hello I'm a client!!\n")

	conn,err:= grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("could not connect to server %v", err)
	}

	defer  conn.Close()

	c:= greetpb.NewGreetServiceClient(conn)
	fmt.Printf("Reponse: %f", c)
}
