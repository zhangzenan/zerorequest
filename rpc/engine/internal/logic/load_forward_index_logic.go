package logic

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"zerorequest/rpc/engine/internal/common/factory"

	"zerorequest/rpc/engine/internal/common/model"
	"zerorequest/rpc/engine/internal/svc"
	"zerorequest/rpc/engine/proto/pb"

	"github.com/edsrzf/mmap-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoadForwardIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoadForwardIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoadForwardIndexLogic {
	return &LoadForwardIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoadForwardIndexLogic) LoadForwardIndex(in *pb.DumpMsg) (*pb.Response, error) {
	// todo: add your logic here and delete this line
	forwardIndex, err := LoadForwardIndex(in.DumpPath)
	if err != nil {
		return nil, err
	}

	return &pb.Response{
		Ok:  true,
		Msg: fmt.Sprintf("加载条数:%d", forwardIndex.Count),
	}, nil
}

const (
	forwardMagic = 0x12345678
)

func LoadForwardIndex(path string) (*model.ForwardIndex, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	mm, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		return nil, err
	}

	idx := &model.ForwardIndex{
		Data:   mm,
		Offset: make(map[uint64]uint64, 1024*1024),
	}

	if err := parseAndBuildIndex(idx); err != nil {
		mm.Unmap()
		return nil, err
	}

	manager := factory.GetForwardIndexManager()

	if err := manager.LoadForwardIndex("product_forward", idx); err != nil {
		mm.Unmap()
		return nil, err
	}

	return idx, nil
}

func parseAndBuildIndex(idx *model.ForwardIndex) error {
	buf := idx.Data
	pos := 0

	magic := binary.LittleEndian.Uint32(buf[pos:])
	pos += 4
	if magic != forwardMagic {
		return fmt.Errorf("invalid magic number: %d", magic)
	}

	version := binary.LittleEndian.Uint32(buf[pos:])
	pos += 4
	_ = version //可校验

	count := binary.LittleEndian.Uint32((buf[pos:]))
	pos += 4
	idx.Count = count

	//顺序扫描所有record,构建offset
	for i := uint32(0); i < count; i++ {
		recOffset := uint64(pos)

		productID := binary.LittleEndian.Uint64(buf[pos:])
		pos += 8

		// 跳过 fixed fields
		pos += 1 // Status
		pos += 4 // Category
		pos += 4 // Stock
		pos += 4 // Price
		pos += 4 // Flags

		tagsLen := binary.LittleEndian.Uint16(buf[pos:])
		pos += 2
		pos += int(tagsLen) // 跳过 tags

		// 建索引
		idx.Offset[productID] = recOffset
	}
	return nil
}
