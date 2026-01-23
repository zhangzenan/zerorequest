package main

import (
	"log"
	"net"
	user "zerorequest/rpc/proto"

	"zerorequest/rpc/user-server/server"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, server.NewUserServer())
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
