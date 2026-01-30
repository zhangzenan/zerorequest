package factory

/*
import (
	"sync"
	"zerorequest/rpc/engine/internal/common/model"
)

// IndexType 索引类型枚举
type IndexType int

const (
	ForwardIndexType  IndexType = iota // 正排索引
	InvertedIndexType                  // 倒排索引
)

// 定义索引接口约束
type Index interface {
	*model.InvertedIndex | *model.ForwardIndex // 支持多种索引类型
}

// 多索引管理器
type IndexManager[T Index] struct {
	indexes map[string]T // 按名称存储多个索引
	mutex   sync.RWMutex
	loading map[string]*sync.Mutex // 正在加载的索引锁
}

// 全局单例实例
var forwardIndexManager = &IndexManager[*model.ForwardIndex]{
	indexes: make(map[string]*model.ForwardIndex),
	loading: make(map[string]*sync.Mutex),
}
var invertedIndexManager = &IndexManager[*model.InvertedIndex]{
	indexes: make(map[string]*model.InvertedIndex),
	loading: make(map[string]*sync.Mutex),
}

func GetIndexManager[T Index](indexTye IndexType) *IndexManager[T] {
	if indexTye == ForwardIndexType {
		return forwardIndexManager
	} else if indexTye == InvertedIndexType {
		return invertedIndexManager
	}
	return nil
}

func loadMutexForIndex[T Index](fm *IndexManager[T], name string) *sync.Mutex {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	if _, exists := fm.loading[name]; !exists {
		fm.loading[name] = &sync.Mutex{}
	}
	return fm.loading[name]
}

func LoadInvertedIndex[T Index](fm *IndexManager[T], name string, index T) error {
	// 为特定索引获取加载锁
	loadMutex := loadMutexForIndex(fm, name)
	loadMutex.Lock()
	defer loadMutex.Unlock()

	// 检查索引是否已经加载
	fm.mutex.RLock()
	if _, exists := fm.indexes[name]; exists {
		fm.mutex.RUnlock()
		return nil
	}
	fm.mutex.RUnlock()

	//存储索引
	fm.mutex.Lock()
	fm.indexes[name] = index
	fm.mutex.Unlock()

	return nil
}

// GetIndexByName 按名称获取索引
func (fm *InvertedIndexManager) GetIndexByName(name string) (*model.InvertedIndex, bool) {
	fm.mutex.RLock()
	defer fm.mutex.RUnlock()

	index, exists := fm.indexes[name]
	return index, exists
}
*/
