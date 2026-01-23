package server

import (
	"context"
	"fmt"
	user "zerorequest/rpc/user/proto/pb"
)

func NewUserServer() user.UserServiceServer {
	return &userServiceServer{}
}

type userServiceServer struct {
	user.UnimplementedUserServiceServer
}

func (s *userServiceServer) AddUser(ctx context.Context, in *user.UserMsg) (*user.Response, error) {
	fmt.Printf("AddUser: %v\n", in)
	return &user.Response{
		Ok: true,
	}, nil
}
