package pkg

import (
	"strconv"
	"strings"
)

func ParseU64(s string) uint64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		// 生产可加监控/日志
		return 0
	}
	return v
}

func ParseU32(s string) uint32 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(v)
}

func ParseU8(s string) uint8 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0
	}
	return uint8(v)
}
