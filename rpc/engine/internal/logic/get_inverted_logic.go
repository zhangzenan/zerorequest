package logic

import (
	"context"
	"unsafe"
	"zerorequest/rpc/engine/internal/common/factory"
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
	postingListView, _ := GetPosting(in.ProductId)
	return &pb.InvertedResponse{
		ProductIds: postingListView.IDs,
	}, nil
}

type PostingListView struct {
	Trigger uint32
	IDs     []uint32
}

func GetPosting(trigger uint32) (PostingListView, bool) {
	manager := factory.GetInvertedIndexManager()
	idx, exists := manager.GetIndexByName("product_inverted")
	if !exists {
		return PostingListView{}, false
	}
	off := idx.Offset[trigger]
	if off == 0 {
		return PostingListView{}, false
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

	return PostingListView{
		Trigger: trigger,
		IDs:     ids,
	}, true

}
