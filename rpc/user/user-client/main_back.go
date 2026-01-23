package main

import (
	"context"
	"fmt"
	"log"
	"time"
	user "zerorequest/rpc/user/proto/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main1() {
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := user.NewUserServiceClient(conn)
	ctx := context.Background()
	in := &user.UserMsg{
		Name:     "张三",
		Age:      18,
		Sex:      user.Sex_FEMALE,
		Birthday: timestamppb.New(time.Now()),
		Hobby:    []string{"羽毛球", "篮球"},
	}
	response, error := client.AddUser(ctx, in)
	if error != nil {
		log.Fatal(error)
	}
	fmt.Println(response)
}
