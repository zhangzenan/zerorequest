package server

import (
	"context"
	"fmt"
	"testing"
	"zerorequest/pkg"
	"zerorequest/rpc/engine/internal/logic"
	"zerorequest/rpc/engine/internal/svc"
	"zerorequest/rpc/engine/proto/pb"
)

var dumpPath = "/data/forward/forward.dump"
var invertedDumpPath = "/data/inverted/inverted.dump"

func TestDataEngineServer_LoadForwardIndex(t *testing.T) {
	// 创建 mock 服务上下文
	svcCtx := svc.NewServiceContext(pkg.Config{})

	load_index_logic := logic.NewLoadForwardIndexLogic(context.Background(), svcCtx)
	in := &pb.DumpMsg{
		DumpPath: dumpPath,
	}
	response, err := load_index_logic.LoadForwardIndex(in)
	if err != nil {
		t.Errorf("LoadForwardIndex() error = %v", err)
		return
	}
	fmt.Printf("LoadForwardIndex() response = %v", response)

	get_forward_logic := logic.NewGetForwardLogic(context.Background(), svcCtx)
	filter := buildFilter()
	forward_request := &pb.ForwardRequest{
		ProductId: 100,
		Filter:    filter,
	}
	forwardResponse, err := get_forward_logic.GetForward(forward_request)
	if err != nil {
		t.Errorf("GetForward() error = %v", err)
		return
	}
	fmt.Printf("GetForward() response = %v", forwardResponse)
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
func TestDataEngineServer_LoadInvertedIndex(t *testing.T) {
	svcCtx := svc.NewServiceContext(pkg.Config{})
	load_forward_logic := logic.NewLoadForwardIndexLogic(context.Background(), svcCtx)
	load_forward_logic.LoadForwardIndex(&pb.DumpMsg{DumpPath: dumpPath})
	load_index_logic := logic.NewLoadInvertedIndexLogic(context.Background(), svcCtx)
	response, err := load_index_logic.LoadInvertedIndex(&pb.DumpMsg{DumpPath: invertedDumpPath})
	if err != nil {
		t.Errorf("LoadInvertedIndex() error = %v", err)
		return
	}
	fmt.Printf("LoadInvertedIndex() response = %v", response)

	get_inverted_logic := logic.NewGetInvertedLogic(context.Background(), svcCtx)

	filter := buildFilter()
	invertedResponse, err := get_inverted_logic.GetInverted(&pb.InvertedRequest{
		ProductIds: []uint32{10, 15},
		Filter:     filter,
		Limit:      50,
	})
	if err != nil {
		t.Errorf("GetInverted() error = %v", err)
		return
	}
	fmt.Printf("GetInverted() response = %v", len(invertedResponse.Results))

}
