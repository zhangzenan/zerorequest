package logic

import (
	"context"
	"encoding/binary"
	"zerorequest/rpc/engine/internal/common/factory"
	"zerorequest/rpc/engine/internal/common/model"
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

	manager := factory.GetForwardIndexManager()
	idx, _ := manager.GetIndexByName("product_forward")
	forwardRecordView, _ := GetProductForward(idx, in.ProductId)

	conditions := make([]Condition, len(in.Filter.Conditions))
	for i, pbCond := range in.Filter.Conditions {
		conditions[i] = Condition{
			Field:  pbCond.Field,
			Op:     Op(pbCond.Op),
			Values: pbCond.Values,
		}
	}
	filter := &Filter{Conds: conditions}
	result := filter.Match(forwardRecordView)
	if !result {
		return nil, nil
	}
	return &pb.ForwardResponse{
		Data: &pb.ForwardView{
			ProductId: forwardRecordView.ProductID,
			Status:    uint32(forwardRecordView.Status),
			Category:  forwardRecordView.Category,
			Stock:     forwardRecordView.Stock,
			Price:     forwardRecordView.Price,
			Flags:     forwardRecordView.Flags,
			Tags:      string(forwardRecordView.Tags),
		},
	}, nil
}

type ForwardRecordView struct {
	ProductID uint32
	Status    uint8
	Category  uint32
	Stock     uint32
	Price     uint32
	Flags     uint32
	Tags      []byte // 指向 mmap 内存
}

func GetProductForward(idx *model.ForwardIndex, productID uint32) (*ForwardRecordView, bool) {
	off, ok := idx.Offset[productID]
	if !ok {
		return nil, false
	}

	buf := idx.Data
	pos := int(off)

	pid := binary.LittleEndian.Uint32(buf[pos:])
	pos += 4

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

type Op uint8

const (
	OpEq Op = iota
	OpNotEq
	OpIn
	OpNotIn
	OpRange
)

type FieldGetter func(*ForwardRecordView) uint64

var FieldGetters = map[string]FieldGetter{
	"price": func(r *ForwardRecordView) uint64 {
		return uint64(r.Price)
	},
	"category": func(r *ForwardRecordView) uint64 {
		return uint64(r.Category)
	},
	"status": func(r *ForwardRecordView) uint64 {
		return uint64(r.Status)
	},
}

type Condition struct {
	Field  string
	Op     Op
	Values []uint64
}

// 单条件执行器
func evalCond(r *ForwardRecordView, c Condition) bool {
	getter := FieldGetters[c.Field]
	if getter == nil {
		return false
	}
	v := getter(r)
	switch c.Op {
	case OpEq:
		return v == c.Values[0]
	case OpNotEq:
		return v != c.Values[0]
	case OpIn:
		for _, x := range c.Values {
			if v == x {
				return true
			}
		}
		return false
	case OpNotIn:
		for _, x := range c.Values {
			if v == x {
				return false
			}
		}
		return true

	case OpRange:
		return v >= c.Values[0] && v <= c.Values[1]
	}
	return false
}

// 条件组合（AND）
type Filter struct {
	Conds []Condition
}

// 默认And
func (f *Filter) Match(r *ForwardRecordView) bool {
	for _, c := range f.Conds {
		if !evalCond(r, c) {
			return false
		}
	}
	return true
}
