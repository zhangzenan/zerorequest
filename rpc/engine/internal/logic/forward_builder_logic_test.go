package logic

import (
	"context"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
	"zerorequest/pkg"
	"zerorequest/rpc/engine/internal/svc"
	"zerorequest/rpc/engine/proto/pb"

	"github.com/stretchr/testify/assert"
)

var csvPath = "D:\\data\\forward\\forward.csv"
var dumpPath = "D:\\data\\forward\\forward.dump"

func TestGenerateProductCSV(t *testing.T) {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 定义表头
	headers := []string{
		"product_id",
		"status",
		"category_id",
		"stock",
		"price",
		"flags",
		"tags",
	}

	// 创建或打开文件
	file, err := os.Create(csvPath)
	if err != nil {
		fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 创建 CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	if err := writer.Write(headers); err != nil {
		fmt.Errorf("写入表头失败: %v", err)
	}

	// 预定义标签列表
	tagsList := []string{"electronics", "clothing", "home", "sports", "books", "toys", "beauty", "health"}

	// 生成商品数据
	numProducts := 1000
	for i := 1; i <= numProducts; i++ {
		product := map[string]string{
			"product_id":  strconv.Itoa(i),
			"status":      strconv.Itoa(rand.Intn(2)),              // 随机选择 0 或 1
			"category_id": strconv.Itoa(rand.Intn(20) + 1),         // 随机数 1-20
			"stock":       strconv.Itoa(rand.Intn(1001)),           // 随机数 0-1000
			"price":       strconv.Itoa(rand.Intn(9999) + 1),       // 随机数 1-10000
			"flags":       strconv.Itoa([]int{0, 3}[rand.Intn(2)]), // 随机选择 0 或 3
			"tags":        getRandomTags(tagsList),                 // 随机生成标签
		}

		// 将 map 转换为切片以便写入 CSV
		record := []string{
			product["product_id"],
			product["status"],
			product["category_id"],
			product["stock"],
			product["price"],
			product["flags"],
			product["tags"],
		}

		// 写入一行数据
		if err := writer.Write(record); err != nil {
			fmt.Errorf("写入数据失败: %v", err)
		}
	}
}

// 获取随机标签
func getRandomTags(tagsList []string) string {
	numTags := rand.Intn(5) + 1 // 随机选择 1-5 个标签
	chosenTags := make([]string, numTags)

	// 随机选择标签
	for i := 0; i < numTags; i++ {
		randomIndex := rand.Intn(len(tagsList))
		chosenTags[i] = tagsList[randomIndex]
	}

	return strings.Join(chosenTags, ",")
}

func TestForwardBuilder(t *testing.T) {
	// 创建 mock 服务上下文
	svcCtx := svc.NewServiceContext(pkg.Config{})

	// 准备测试数据
	in := &pb.DumpMsg{
		CsvPath:  csvPath,
		DumpPath: dumpPath,
	}

	logic := NewForwardBuilderLogic(context.Background(), svcCtx)
	resp, err := logic.ForwardBuilder(in)

	assert.NoError(t, err)
	assert.True(t, resp.Ok)
}

func TestForwardInspect(t *testing.T) {
	f, err := os.Open(dumpPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var header ForwardHeader
	if err := binary.Read(f, binary.LittleEndian, &header); err != nil {
		panic(err)
	}

	fmt.Printf("Header: %+v\n", header)

	// 只打印前 N 条，防止刷屏
	const maxPrint = 10

	for i := 0; i < int(header.Count) && i < maxPrint; i++ {
		var productID uint64
		var status uint8
		var category uint32
		var stock uint32
		var price uint32
		var flags uint32
		var tagsLen uint16

		binary.Read(f, binary.LittleEndian, &productID)
		binary.Read(f, binary.LittleEndian, &status)
		binary.Read(f, binary.LittleEndian, &category)
		binary.Read(f, binary.LittleEndian, &stock)
		binary.Read(f, binary.LittleEndian, &price)
		binary.Read(f, binary.LittleEndian, &flags)
		binary.Read(f, binary.LittleEndian, &tagsLen)

		tags := make([]byte, tagsLen)
		if _, err := io.ReadFull(f, tags); err != nil {
			panic(err)
		}

		fmt.Printf(
			"[%d] productID=%d status=%d category=%d stock=%d price=%d flags=%d tags=%s\n",
			i, productID, status, category, stock, price, flags, string(tags),
		)
	}
}
