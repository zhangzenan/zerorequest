package main

import (
	"context"
	"fmt"
	"log"
	"time"
	user "zerorequest/rpc/proto"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	// 从 Nacos 获取服务实例
	instance := discoverService()
	if instance == nil {
		log.Fatal("未能发现可用的服务实例")
	}
	// 构建服务地址
	address := fmt.Sprintf("%s:%d", instance.Ip, instance.Port)
	fmt.Printf("连接到服务实例: %s\n", address)

	// 使用服务实例地址建立 gRPC 连接
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

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

func discoverService() *model.Instance {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("nacos-test.com", 8848, constant.WithContextPath("/nacos")),
	}
	cc := constant.ClientConfig{
		NamespaceId: "dev",
		TimeoutMs:   5000,
	}

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		log.Fatalf("创建 Nacos 客户端失败: %v", err)
	}

	// SelectOneHealthyInstance将会按加权随机轮询的负载均衡策略返回一个健康的实例
	// 实例必须满足的条件：health=true,enable=true and weight>0
	instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: "user-service",
		GroupName:   "DEFAULT_GROUP",     // 默认值DEFAULT_GROUP
		Clusters:    []string{"DEFAULT"}, // 默认值DEFAULT
	})

	return instance
}
