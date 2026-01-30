package model

type ForwardIndex struct {
	Data   []byte // mmap 映射的整个文件
	Count  uint32
	Offset map[uint64]uint64 // productID -> offset
}
