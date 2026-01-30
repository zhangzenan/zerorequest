package logic

import (
	"context"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"zerorequest/rpc/engine/internal/svc"
	"zerorequest/rpc/engine/proto/pb"

	"github.com/zeromicro/go-zero/core/logx"

	"zerorequest/pkg"
)

type ForwardBuilderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewForwardBuilderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForwardBuilderLogic {
	return &ForwardBuilderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ForwardBuilderLogic) ForwardBuilder(in *pb.DumpMsg) (*pb.Response, error) {
	// todo: add your logic here and delete this line
	err := buildForward(in.CsvPath, in.DumpPath)
	return &pb.Response{
		Ok: true,
	}, err
}

type ForwardHeader struct {
	Magic   uint32
	Version uint32
	Count   uint32
}

type ForwardRecord struct {
	ProductID uint64
	Status    uint8
	Category  uint32
	Stock     uint32
	Price     uint32
	Flags     uint32
	Tags      []byte
}

func buildForward(csvFile, outFile string) error {
	f, _ := os.Create(outFile)
	defer f.Close()

	//写Header占位
	header := ForwardHeader{
		Magic:   0x12345678,
		Version: 1,
		Count:   0,
	}
	binary.Write(f, binary.LittleEndian, header)
	file, _ := openCSVFile(csvFile)
	reader := csv.NewReader(file)

	// 读取 title
	title, err := reader.Read()
	if err != nil {
		return fmt.Errorf("read csv header failed: %w", err)
	}

	// 可选：校验字段顺序，防止 schema 变更
	if len(title) < 7 ||
		title[0] != "product_id" ||
		title[1] != "status" ||
		title[2] != "category_id" ||
		title[3] != "stock" ||
		title[4] != "price" ||
		title[5] != "flags" ||
		title[6] != "tags" {
		return fmt.Errorf("unexpected csv header: %+v", header)
	}

	var count uint32

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		rec := ForwardRecord{
			ProductID: pkg.ParseU64(row[0]),
			Status:    pkg.ParseU8(row[1]),
			Category:  pkg.ParseU32(row[2]),
			Stock:     pkg.ParseU32(row[3]),
			Price:     pkg.ParseU32(row[4]),
			Flags:     pkg.ParseU32(row[5]),
			Tags:      []byte(row[6]),
		}
		// 写 fixed
		binary.Write(f, binary.LittleEndian, rec.ProductID)
		binary.Write(f, binary.LittleEndian, rec.Status)
		binary.Write(f, binary.LittleEndian, rec.Category)
		binary.Write(f, binary.LittleEndian, rec.Stock)
		binary.Write(f, binary.LittleEndian, rec.Price)
		binary.Write(f, binary.LittleEndian, rec.Flags)
		// 写变长
		tagsLen := uint16(len(rec.Tags))
		binary.Write(f, binary.LittleEndian, tagsLen)
		f.Write(rec.Tags)

		count++
	}

	//回写Header Count
	f.Seek(0, 0)
	header.Count = count
	binary.Write(f, binary.LittleEndian, header)
	return nil
}

func openCSVFile(path string) (*os.File, error) {
	return os.Open(path)
}
