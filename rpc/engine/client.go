package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
	"zerorequest/rpc/engine/proto/pb"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	client := pb.NewDataEngineClient(conn)
	ctx := context.Background()

	request(client, ctx)
}

func discoverService() *model.Instance {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("nacos-test.shwoody.com", 8848, constant.WithContextPath("/nacos")),
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
		ServiceName: "data-engine-service",
		GroupName:   "DEFAULT_GROUP",     // 默认值DEFAULT_GROUP
		Clusters:    []string{"DEFAULT"}, // 默认值DEFAULT
	})

	return instance
}

func request(client pb.DataEngineClient, ctx context.Context) {

	filter := buildFilter()
	for i := 0; i < 10000000; i++ {
		var wg sync.WaitGroup

		//每次启动50个并行请求
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// 生成200个1到1000之间的随机数切片
				productIds := generateRandomProductIds(200, 1, 10000000)
				response, error := client.GetInverted(ctx, &pb.InvertedRequest{
					ProductIds: productIds,
					Filter:     filter,
					Limit:      50,
				})
				if error != nil {
					fmt.Printf("请求失败: %v", error)
				}
				if rand.Float64() < 0.1 {
					fmt.Printf("index:%d 请求结果length: %v\n", i, len(response.Results))
				}
			}()
		}
		wg.Wait() //等待当前批次的请求结束

		time.Sleep(500 * time.Millisecond)
	}
}

// 生成指定数量的随机产品ID
func generateRandomProductIds(count, min, max int) []uint32 {
	productIds := make([]uint32, count)
	for i := 0; i < count; i++ {
		productIds[i] = uint32(rand.Intn(max-min+1) + min)
	}
	return productIds
}

func buildFilter() *pb.Filter {
	filter := &pb.Filter{
		Conditions: []*pb.Condition{
			{
				Field:  "status",
				Op:     pb.Operation_OpEq,
				Values: []uint64{1},
			},
			{
				Field:  "price",
				Op:     pb.Operation_OpRange,
				Values: []uint64{20, 1000},
			},
		},
	}
	return filter
}
