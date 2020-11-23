package main

import (
	"context"
	"fmt"
	pb "mygrpc/messages"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)


const (
	port = ":50051"
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
	socket, err :=	net.Listen("tcp", port)

	if err != nil{
		fmt.Println(err)
	}
	
	fmt.Println("Starting grpc server at "+port)
	s := grpc.NewServer()
	pb.RegisterAgentServer(s, &server{})
	s.Serve(socket)
}


