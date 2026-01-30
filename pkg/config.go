package pkg

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	CusLog LogConf
}

type LogConf struct {
	Level       string `json:"level"`
	Filename    string `json:"filename"`
	MaxSize     int    `json:"maxSize"` //MB
	MaxAge      int    `json:"maxAge"`  //天
	MaxBackups  int    `json:"maxBackups"`
	Compress    bool   `json:"compress"`
	ServiceName string `json:"serviceName"`
}
