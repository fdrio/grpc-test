package main

import (
	"context"
	"fmt"
	pb "mygrpc/messages"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)



type server struct{
	pb.UnimplementedAgentServer
}

func (s *server) EchoReply(ctx context.Context, reply *pb.Reply)(*pb.Reply, error){
	fmt.Println("Server got: "+reply.GetInfo())

	newReply := &pb.Reply{
		Timestamp:timestamppb.Now(),
		Info: "Replying back: "+reply.GetInfo(),
}
	return newReply, nil
}

func main(){
	arguments := os.Args

	if len(arguments)==1{
		fmt.Println("Invalid: Enter port number")
	}
	
	port := ":"+arguments[1]
	socket, err :=	net.Listen("tcp", port)
	if err != nil{
		fmt.Println(err)
	}
	
	defer socket.Close()
	fmt.Println("Starting grpc server at "+port)
	s := grpc.NewServer()
	pb.RegisterAgentServer(s, &server{})
	s.Serve(socket)
}


