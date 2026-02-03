package logic

import (
	"context"
	"math/rand"
	"sync"
	"time"
	"unsafe"
	"zerorequest/rpc/engine/internal/common/factory"
	"zerorequest/rpc/engine/internal/common/model"
	"zerorequest/rpc/engine/internal/svc"
	"zerorequest/rpc/engine/proto/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInvertedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetInvertedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInvertedLogic {
	return &GetInvertedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetInvertedLogic) GetInverted(in *pb.InvertedRequest) (*pb.InvertedResponse, error) {
	// todo: add your logic here and delete this line
	startTime := time.Now()
	//postingListView, _ := GetPosting(in.ProductIds[0])
	forwardRecordView := BatchGetPosting(in)
	duration := time.Since(startTime)
	// 计算微秒和毫秒
	durationMicroseconds := duration.Microseconds()
	durationMilliseconds := duration.Seconds() * 1000

	// 10% 采样打印日志
	if rand.Float64() < 0.1 {
		logx.Infof("GetPosting 耗时: %d 微秒, %.2f 毫秒", durationMicroseconds, durationMilliseconds)
	}

	result := make(map[uint32]*pb.ForwardViewList, len(in.ProductIds))
	for key, forwardRecordViewList := range forwardRecordView {
		// 创建 ForwardView 切片
		forwardViews := make([]*pb.ForwardView, 0, len(forwardRecordViewList))

		for _, forwardRecordView := range forwardRecordViewList {
			forwardViews = append(forwardViews, &pb.ForwardView{
				ProductId: forwardRecordView.ProductID,
				Status:    uint32(forwardRecordView.Status),
				Category:  forwardRecordView.Category,
				Tags:      string(forwardRecordView.Tags),
				Price:     forwardRecordView.Price,
				Stock:     forwardRecordView.Stock,
				Flags:     forwardRecordView.Flags,
			})
		}
		// 创建 ForwardViewList 并设置 Data 字段
		result[key] = &pb.ForwardViewList{
			Data: forwardViews,
		}
	}
	return &pb.InvertedResponse{
		Results: result,
	}, nil
}

type PostingListView struct {
	Trigger uint32
	IDs     []uint32
}

func GetPosting(idx *model.InvertedIndex, trigger uint32) (*PostingListView, bool) {
	off := idx.Offset[trigger]
	if off == 0 {
		return &PostingListView{}, false
	}
	// 使用 unsafe.Pointer 进行指针运算
	/**
	idx.Data[off]：获取数据数组中偏移量 off 位置的单个字节
	&idx.Data[off]：获取该字节的内存地址（*byte 类型）
	unsafe.Pointer(...)：将 *byte 类型转换为通用指针类型 unsafe.Pointer
	*/
	cntPtr := unsafe.Pointer(&idx.Data[off])

	//cnt
	cnt := *(*uint32)(cntPtr)

	//rid array
	// 计算 IDs 数组的起始位置 (跳过 4 字节 cnt)
	idsPtr := unsafe.Pointer(uintptr(cntPtr) + 4)
	ids := unsafe.Slice((*uint32)(idsPtr), cnt)

	return &PostingListView{
		Trigger: trigger,
		IDs:     ids,
	}, true
}

func BatchGetPosting(in *pb.InvertedRequest) map[uint32][]*ForwardRecordView {
	conditions := make([]Condition, len(in.Filter.Conditions))
	for i, pbCond := range in.Filter.Conditions {
		conditions[i] = Condition{
			Field:  pbCond.Field,
			Op:     Op(pbCond.Op),
			Values: pbCond.Values,
		}
	}
	filter := &Filter{Conds: conditions}
	ids := in.ProductIds
	const groupSize = 20
	// 将IDs分组
	var groups [][]uint32
	for i := 0; i < len(ids); i += groupSize {
		end := i + groupSize
		if end > len(ids) {
			end = len(ids)
		}
		groups = append(groups, ids[i:end])
	}
	manager := factory.GetInvertedIndexManager()
	inverted_idx, _ := manager.GetIndexByName("product_inverted")
	forward_manager := factory.GetForwardIndexManager()
	forward_idx, _ := forward_manager.GetIndexByName("product_forward")

	var finalResult sync.Map
	var wg sync.WaitGroup // 声明 WaitGroup
	// 并发执行每个组
	for _, group := range groups {
		wg.Add(1) // 每启动一个 Goroutine，计数器加一
		go func(g []uint32) {
			defer wg.Done() // Goroutine 结束时调用 Done 减少计数器

			groupResult := make(map[uint32][]*ForwardRecordView)

			// 在组内串行执行GetPosting
			for _, id := range g {
				postingList, exists := GetPosting(inverted_idx, id)
				if !exists {
					continue
				}
				limit := uint32(0)
				for _, rid := range postingList.IDs {
					if limit >= in.Limit {
						break
					}
					forwardRecordView, _ := GetProductForward(forward_idx, rid)
					if forwardRecordView == nil {
						continue
					}
					filter_result := filter.Match(forwardRecordView)
					if filter_result {
						groupResult[id] = append(groupResult[id], forwardRecordView)
						limit += 1
					}
				}
			}
			for k, v := range groupResult {
				finalResult.Store(k, v)
			}
		}(group)
	}
	wg.Wait() // 等待所有 Goroutine 完成
	// 收集所有组的结果
	// 最后统一读取 sync.Map 中的数据
	result := make(map[uint32][]*ForwardRecordView)
	finalResult.Range(func(key, value interface{}) bool {
		result[key.(uint32)] = value.([]*ForwardRecordView)
		return true
	})

	return result
}
