package model

type InvertedHeader struct {
	Magic    uint32
	Version  uint32
	KeyCount uint32
}

type KeyIndexEntry struct {
	TriggerID uint32
	Offset    uint64 // PostingList 在文件中的偏移
}

type InvertedIndex struct {
	Data   []byte // mmap 全文件
	Count  uint32
	Offset map[uint32]uint64
}
