package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"zerorequest/pkg"
	"zerorequest/rpc/user/internal/server"
	"zerorequest/rpc/user/internal/svc"
	user "zerorequest/rpc/user/proto/pb"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/zeromicro/go-zero/core/conf"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "../etc/user.yaml", "the config file")

func main() {
	//配置Nacos服务端
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("nacos-test.shwoody.com", 8848, constant.WithContextPath("/nacos")),
	}

	//配置Nacos客户端
	cc := constant.ClientConfig{
		NamespaceId:         "dev",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	//创建服务发现客户端
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: sc,
			ClientConfig:  &cc,
		},
	)
	if err != nil {
		log.Fatalf("创建 Nacos 客户端失败: %v", err)
	}
	//启动gRPC 服务
	lis, err := net.Listen("tcp", ":0") //使用：0自动分配端口
	if err != nil {
		log.Fatalf("启动gRPC服务失败: %v", err)
	}
	port := lis.Addr().(*net.TCPAddr).Port //获取实际分配的端口
	fmt.Printf("gRPC服务已启动，端口为: %d\n", port)

	s := grpc.NewServer()
	var c pkg.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	user.RegisterUserServiceServer(s, server.NewUserServiceServer(ctx))

	//在goroutine中启动gRPC服务
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("启动gRPC服务失败: %v", err)
		}
	}()

	//向Nacos注册服务
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          getLocalIP(),   // 服务实例 IP
		Port:        uint64(port),   // 服务实例端口
		ServiceName: "user-service", // 服务名称
		Weight:      10,             // 权重
		Enable:      true,           // 是否启用
		Healthy:     true,           // 健康状态
		Metadata: map[string]string{
			"protocol": "grpc",  // 元数据 - 协议
			"version":  "1.0.0", // 元数据 - 版本
		},
		ClusterName: "DEFAULT",       // 集群名称
		GroupName:   "DEFAULT_GROUP", // 分组名称
		Ephemeral:   true,            // 是否临时实例
	})

	if err != nil {
		log.Fatalf("注册服务到 Nacos 失败: %v", err)
	}

	if !success {
		log.Fatal("服务注册未成功")
	}

	fmt.Println("服务已成功注册到 Nacos")
	// 优雅关闭处理
	defer func() {
		// 从 Nacos 注销服务
		_, _ = namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          getLocalIP(),
			Port:        uint64(port),
			ServiceName: "user-service",
			GroupName:   "DEFAULT_GROUP",
		})
		s.GracefulStop()
	}()

	// 保持服务运行
	select {}

}

// 获取机器的真实IP地址
func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
