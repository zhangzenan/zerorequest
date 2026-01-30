package logic

import (
	"context"
	"encoding/binary"
	"zerorequest/rpc/engine/internal/common/factory"
	"zerorequest/rpc/engine/internal/svc"
	"zerorequest/rpc/engine/proto/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetForwardLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetForwardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetForwardLogic {
	return &GetForwardLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetForwardLogic) GetForward(in *pb.ForwardRequest) (*pb.ForwardResponse, error) {
	// todo: add your logic here and delete this line

	forwardRecordView, _ := GetProductForward(in.ProductId)
	return &pb.ForwardResponse{
		ProductId: forwardRecordView.ProductID,
		Status:    uint32(forwardRecordView.Status),
		Category:  forwardRecordView.Category,
		Stock:     forwardRecordView.Stock,
		Price:     forwardRecordView.Price,
		Flags:     forwardRecordView.Flags,
		Tags:      string(forwardRecordView.Tags),
	}, nil
}

type ForwardRecordView struct {
	ProductID uint64
	Status    uint8
	Category  uint32
	Stock     uint32
	Price     uint32
	Flags     uint32
	Tags      []byte // 指向 mmap 内存
}

func GetProductForward(productID uint64) (*ForwardRecordView, bool) {
	manager := factory.GetForwardIndexManager()
	idx, exists := manager.GetIndexByName("product_forward")
	if !exists {
		return nil, false
	}
	off, ok := idx.Offset[productID]
	if !ok {
		return nil, false
	}

	buf := idx.Data
	pos := int(off)

	pid := binary.LittleEndian.Uint64(buf[pos:])
	pos += 8

	status := buf[pos]
	pos += 1

	category := binary.LittleEndian.Uint32(buf[pos:])
	pos += 4

	stock := binary.LittleEndian.Uint32(buf[pos:])
	pos += 4

	price := binary.LittleEndian.Uint32(buf[pos:])
	pos += 4

	flags := binary.LittleEndian.Uint32(buf[pos:])
	pos += 4

	tagsLen := binary.LittleEndian.Uint16(buf[pos:])
	pos += 2

	// 优化：创建 tags 的副本，避免内存共享
	tags := make([]byte, tagsLen)
	copy(tags, buf[pos:pos+int(tagsLen)])

	return &ForwardRecordView{
		ProductID: pid,
		Status:    status,
		Category:  category,
		Stock:     stock,
		Price:     price,
		Flags:     flags,
		Tags:      tags,
	}, true
}
