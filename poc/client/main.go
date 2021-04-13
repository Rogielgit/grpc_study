package main

import (
	"context"
	"example.com/grpc_study/poc/poc_proto"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
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

	fmt.Println("Performing client requests")

	rand.Seed(time.Now().UnixNano())
	clientID := rand.Intn(1000)
	fmt.Printf("clientID %d\n", clientID)

	for {
		res := getWorkResponse(conn, c)
		switch {
		case res.GetCheck().HasWork:
			doWork(c, clientID)
		default:
			fmt.Println("client does not have work to do!!")
			time.Sleep(2000 * time.Millisecond)
		}
	}
}

func isDeadlineExceeded(err error) {
	s, ok := status.FromError(err)
	if ok {
		if s.Code() == codes.DeadlineExceeded {
			log.Printf("deadline exceeded\n", err)

		}
	}
}

func getWorkResponse(conn *grpc.ClientConn, c poc_proto.CheckServiceClient) (res *poc_proto.CheckWorkResponse) {
	res, err := askServerWork(c)
	if err != nil {
		log.Printf("unexpected error from HasWork: %v\n", err)
		for {

			if conn.GetState() == connectivity.Ready {
				res, _ = askServerWork(c)
				break
			} else {
				time.Sleep(2000 * time.Millisecond)
				fmt.Println("trying to reconnect to the server")
			}
		}
	}
	return
}

func askServerWork(c poc_proto.CheckServiceClient) (res *poc_proto.CheckWorkResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second) //Deadline
	defer cancel()
	return c.HasWork(ctx, &poc_proto.CheckWorkRequest{})
}

func doWork(c poc_proto.CheckServiceClient, clientID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second) //Deadline
	defer cancel()
	fmt.Println("client has work to do!!")
	stream, err := c.DoWork(ctx)
	for i := 0; i < 10; i++ {
		network := strconv.Itoa(clientID) + "-network" + strconv.Itoa(i)
		if err != nil {
			log.Fatalf("error while send info to server")
		}
		stream.Send(&poc_proto.DoWorkRequest{
			Network: network,
		})
	}

	endResp, errResp := stream.CloseAndRecv()
	if errResp != nil {
		log.Fatalf("error while receiving response from server", errResp)
	}
	fmt.Printf("server response to client %v: %v\n ", clientID, endResp.GetMessage())
}
