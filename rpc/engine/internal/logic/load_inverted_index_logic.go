package logic

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"unsafe"
	"zerorequest/rpc/engine/internal/common/factory"
	"zerorequest/rpc/engine/internal/common/model"

	"zerorequest/rpc/engine/internal/svc"
	"zerorequest/rpc/engine/proto/pb"

	"github.com/edsrzf/mmap-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoadInvertedIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoadInvertedIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoadInvertedIndexLogic {
	return &LoadInvertedIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoadInvertedIndexLogic) LoadInvertedIndex(in *pb.DumpMsg) (*pb.Response, error) {
	// todo: add your logic here and delete this line
	invertedIndex, err := LoadInvertedIndex(in.DumpPath)
	return &pb.Response{
		Ok:  err != nil,
		Msg: fmt.Sprintf("加载条数:%d", invertedIndex.Count),
	}, nil
}

func LoadInvertedIndex(path string) (*model.InvertedIndex, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	mm, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		return nil, err
	}

	idx := &model.InvertedIndex{
		Data: mm,
	}
	//1.Header
	header := (*model.InvertedHeader)(unsafe.Pointer(&mm[0]))
	if header.Magic != 0x12345678 {
		return nil, fmt.Errorf("invalid magic number: %d", header.Magic)
	}

	//2.KeyIndex
	off := uintptr(unsafe.Sizeof(model.InvertedHeader{})) //记录相对于文件起始位置的字节偏移量

	idx.Count = header.KeyCount
	idx.Offset = make(map[uint32]uint64, idx.Count)
	for i := 0; i < int(idx.Count); i++ {
		triggerID := binary.LittleEndian.Uint32(mm[off:])
		off += 4
		offset := binary.LittleEndian.Uint64(mm[off:])
		off += 8
		//entry := (*model.KeyIndexEntry)(unsafe.Pointer(&mm[off])) //将内存映射中指定偏移量的数据转换为 KeyIndexEntry 结构体指针，结构体大小固定，无需额外的结束标记
		idx.Offset[triggerID] = offset
		//size := unsafe.Sizeof(model.KeyIndexEntry{})
	}

	manager := factory.GetInvertedIndexManager()

	if err := manager.LoadInvertedIndex("product_inverted", idx); err != nil {
		mm.Unmap()
		return nil, err
	}

	return idx, nil
}
