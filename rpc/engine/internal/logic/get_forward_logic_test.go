package logic

import (
	"context"
	"testing"
	"zerorequest/pkg"
	"zerorequest/rpc/engine/internal/svc"
	"zerorequest/rpc/engine/proto/pb"
)

func TestGetForwardLogic_GetForward(t *testing.T) {
	// 创建 mock 服务上下文
	svcCtx := svc.NewServiceContext(pkg.Config{})

	// 准备测试数据
	in := &pb.ForwardRequest{
		ProductId: 1,
	}

	logic := NewGetForwardLogic(context.Background(), svcCtx)
	forwardResponse, err := logic.GetForward(in)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(forwardResponse)
}
