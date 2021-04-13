package main

import (
	"context"
	"example.com/grpc_study/poc/poc_proto"
	"fmt"
	"google.golang.org/grpc/connectivity"
	"log"
	"math/rand"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

func main() {

	conn, err := grpc.Dial("0.0.0.0:5051", grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := poc_proto.NewCheckServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute) //Deadline
	defer cancel()
	fmt.Println("Performing unary request")

	rand.Seed(time.Now().UnixNano())
	clientID := rand.Intn(1000)
	fmt.Printf("clientID %d\n", clientID)

	for {
		res, err := c.HasWork(ctx, &poc_proto.CheckWorkRequest{})
		if err != nil {
			log.Printf("unexpected error from HasWork: %v", err)
			for {
				if conn.GetState() == connectivity.Ready {
					res, _ = c.HasWork(ctx, &poc_proto.CheckWorkRequest{})
					break
				} else {
					time.Sleep(2000 * time.Millisecond)
					fmt.Println("trying to reconnect to the server")
				}
			}
		}
		switch {
		case res.GetCheck().HasWork:
			fmt.Println("client has work to do!!")
			for i := 0; i < 10; i++ {
				network := strconv.Itoa(clientID) + "-network" + strconv.Itoa(i)
				stream, err := c.DoWork(ctx)
				if err != nil {
					log.Fatalf("error while send info to server")
				}
				stream.Send(&poc_proto.DoWorkRequest{
					Network: network,
				})
			}

		default:
			fmt.Println("client does not have work to do!!")
			time.Sleep(2000 * time.Millisecond)
		}
	}
}
