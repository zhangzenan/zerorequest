package logic

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/csv"
	"io"
	"os"
	"sort"
	"strings"
	"unsafe"
	"zerorequest/pkg"

	"zerorequest/rpc/engine/internal/common/model"
	"zerorequest/rpc/engine/internal/svc"
	"zerorequest/rpc/engine/proto/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type InvertedBuilderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInvertedBuilderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InvertedBuilderLogic {
	return &InvertedBuilderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InvertedBuilderLogic) InvertedBuilder(in *pb.DumpMsg) (*pb.Response, error) {
	// todo: add your logic here and delete this line

	err := buildInverted(in.CsvPath, in.DumpPath)
	return &pb.Response{
		Ok: err == nil,
	}, nil
}

type PostingListHeader struct {
	TriggerID uint64
	Count     uint32
	// 后面紧跟 Count 个 uint64 relatedIDs
}

func buildInverted(csvFile, outFile string) error {
	//文件布局，第一部分是header，第二部分是keyIndex，第三部分是PostingLists

	//1.读取csv聚合
	m := make(map[uint32][]uint32)

	file, _ := open(csvFile)
	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		trigger := pkg.ParseU32(row[0])
		// 将逗号分隔的字符串转换为切片
		relatedStr := row[1]
		relatedParts := strings.Split(relatedStr, ",")
		// 遍历所有相关的 ID 并添加到映射中
		for _, part := range relatedParts {
			part = strings.TrimSpace(part) // 去除可能的空白字符
			if part != "" {                // 避免空字符串转换为数字
				related := pkg.ParseU32(part)
				m[trigger] = append(m[trigger], related)
			}
		}
	}
	//2.排序trigger keys
	keys := make([]uint32, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	f, _ := os.Create(outFile)
	defer f.Close()

	// 使用带缓冲的写入器
	writer := bufio.NewWriterSize(f, 1024*1024)

	//写Header
	header := model.InvertedHeader{
		Magic:    0x12345678,
		Version:  1,
		KeyCount: uint32(0), // 初始为 0
	}
	binary.Write(writer, binary.LittleEndian, header)
	writer.Flush() // 立即刷新 header 数据到文件

	//keyIndex占位
	indexPos := make([]int64, len(keys))
	for range keys {
		binary.Write(writer, binary.LittleEndian, uint32(0)) // trigger
		binary.Write(writer, binary.LittleEndian, uint64(0)) // offset
	}
	writer.Flush() // 确保 keyIndex 占位数据写入

	//写PostingLists
	// 在写入前计算预计的文件位置
	//pos := int64(binary.Size(header)) + int64(len(keys)*12) // header大小 + keyIndex区域大小

	for i, k := range keys {
		// 在当前位置写入，而不是使用pos跟踪
		offset, _ := f.Seek(0, io.SeekCurrent) // 获取当前实际位置
		indexPos[i] = offset
		//写PostingList
		cnt := uint32(len(m[k]))
		// 手动编码长度
		lenBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(lenBytes, cnt)
		writer.Write(lenBytes) //PostingList的长度

		// 批量编码 related IDs
		idBytes := make([]byte, len(m[k])*4)
		for j, rid := range m[k] {
			startIdx := j * 4
			binary.LittleEndian.PutUint32(idBytes[startIdx:startIdx+4], rid)
		}
		writer.Write(idBytes)

		//pos += int64(4 + len(m[k])*4) // 更新预计位置：长度(4字节) + ID数量*
	}
	//回写keyIndex
	f.Seek(int64(binary.Size(header)), io.SeekStart) //定位目标：跳转到文件头部之后的第一个字节位置
	for i, k := range keys {
		binary.Write(f, binary.LittleEndian, k)                   //写trigger
		binary.Write(f, binary.LittleEndian, uint64(indexPos[i])) //写offset
	}
	writer.Flush()
	// 在所有数据写入完成后，回写真实的 KeyCount
	realKeyCount := uint32(len(keys))

	// 移动到文件开头的 KeyCount 位置
	f.Seek(int64(unsafe.Sizeof(uint32(0))*2), io.SeekStart) // 跳过 Magic 和 Version
	binary.Write(f, binary.LittleEndian, realKeyCount)

	return nil
}

func open(path string) (*os.File, error) {
	return os.Open(path)
}
