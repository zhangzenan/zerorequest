package logic

import (
	"context"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"zerorequest/pkg"
	"zerorequest/rpc/engine/internal/svc"
	"zerorequest/rpc/engine/proto/pb"

	"github.com/stretchr/testify/assert"
)

var invertedCsvPath = "D:\\data\\inverted\\inverted.csv"
var invertedDumpPath = "D:\\data\\inverted\\inverted.dump"

func TestBuildInvertedCsv(t *testing.T) {
	// 创建测试CSV文件
	file, err := os.Create(invertedCsvPath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 生成 1000 条测试数据
	for i := 0; i < 10000; i++ {
		triggerId := uint64(i + 1)

		// 为每个triggerId生成2-5个相关商品ID
		relatedIds := make([]string, 0, 5)
		relatedCount := rand.Intn(190) + 5 // 随机2-5个相关ID

		for j := 0; j < relatedCount; j++ {
			relatedId := triggerId + uint64(j) // 生成相关的商品ID
			relatedIds = append(relatedIds, fmt.Sprintf("%d", relatedId))
		}

		// 创建CSV行：triggerId, "relatedId1,relatedId2,..."
		row := []string{
			fmt.Sprintf("%d", triggerId),
			strings.Join(relatedIds, ","), // 逗号分隔的相关商品ID
		}

		if err := writer.Write(row); err != nil {
			t.Error(err)
		}
	}
}

func TestInvertedBuilderLogic_InvertedBuilder(t *testing.T) {
	svcCtx := svc.NewServiceContext(pkg.Config{})
	logic := NewInvertedBuilderLogic(context.Background(), svcCtx)
	resp, err := logic.InvertedBuilder(&pb.DumpMsg{
		CsvPath:  invertedCsvPath,
		DumpPath: invertedDumpPath,
	})
	assert.NoError(t, err)
	assert.True(t, resp.Ok)
}
