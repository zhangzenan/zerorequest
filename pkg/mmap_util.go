package pkg

import (
	"os"

	"github.com/edsrzf/mmap-go"
)

func MmapFile(path string) ([]byte, func(), error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	// 使用第三方库的 mmap
	data, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		f.Close()
		return nil, nil, err
	}

	cleanup := func() {
		data.Unmap() // 解除内存映射
		f.Close()
	}

	return data, cleanup, nil
}
