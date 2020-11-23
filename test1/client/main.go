package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "mygrpc/messages"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	address = "localhost:50051"
)

func main(){
	// grpc connection
	// takes address and variadic options 
	conn , err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil{
		// Error connecting
		log.Println(err)
	}

	defer conn.Close()
	client := pb.NewAgentClient(conn)	

	// create a new context
	ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	
	reply := &pb.Reply{Info: "Hello from client", Timestamp:timestamppb.Now() }

	

	// Call grpc server 
	serverReply, err := client.EchoReply(ctx, reply)
	
	if err != nil{
		log.Println("Error calling rpc")
	}

	fmt.Println("Info: "+serverReply.Info)	
	fmt.Println("Timestamp: "+serverReply.Timestamp.AsTime().String())
	defer cancel()

}	
