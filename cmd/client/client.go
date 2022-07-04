package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/sephyret/grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Connect()

	client := pb.NewUserServiceClient(connection)
	// Implementation unaria
	// AddUser(client)

	// Implementation serve stream
	// AddUserVerbose(client)

	// Implementation client stream
	// AddUsers(client)

	// Implementation server stream and client stream
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Leonardo",
		Email: "leo@leo",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not connect to make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Leonardo",
		Email: "leo@leo",
	}

	resStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not connect to make gRPC request: %v", err)
	}

	for {
		stream, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}
		fmt.Println("Status: ", stream.Status)
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:    "leo1",
			Name:  "Leonardo 1",
			Email: "leo1@leo.com",
		},
		{
			Id:    "leo2",
			Name:  "Leonardo 2",
			Email: "leo2@leo.com",
		},
		{
			Id:    "leo3",
			Name:  "Leonardo 3",
			Email: "leo3@leo.com",
		},
		{
			Id:    "leo4",
			Name:  "Leonardo 4",
			Email: "leo4@leo.com",
		},
		{
			Id:    "leo5",
			Name:  "Leonardo 5",
			Email: "leo5@leo.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		{
			Id:    "leo1",
			Name:  "Leonardo 1",
			Email: "leo1@leo.com",
		},
		{
			Id:    "leo2",
			Name:  "Leonardo 2",
			Email: "leo2@leo.com",
		},
		{
			Id:    "leo3",
			Name:  "Leonardo 3",
			Email: "leo3@leo.com",
		},
		{
			Id:    "leo4",
			Name:  "Leonardo 4",
			Email: "leo4@leo.com",
		},
		{
			Id:    "leo5",
			Name:  "Leonardo 5",
			Email: "leo5@leo.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.GetName())
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
			}
			fmt.Printf("Receiving user %v com status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
