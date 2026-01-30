package factory

import (
	"sync"
	"zerorequest/rpc/engine/internal/common/model"
)

// 多索引管理器
type ForwardIndexManager struct {
	indexes map[string]*model.ForwardIndex // 按名称存储多个索引
	mutex   sync.RWMutex
	loading map[string]*sync.Mutex // 正在加载的索引锁
}

// 全局单例实例
var forwardIndexManager = &ForwardIndexManager{
	indexes: make(map[string]*model.ForwardIndex),
	loading: make(map[string]*sync.Mutex),
}

func GetForwardIndexManager() *ForwardIndexManager {
	return forwardIndexManager
}

func (fm *ForwardIndexManager) loadMutexForIndex(name string) *sync.Mutex {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	if _, exists := fm.loading[name]; !exists {
		fm.loading[name] = &sync.Mutex{}
	}
	return fm.loading[name]
}

func (fm *ForwardIndexManager) LoadForwardIndex(name string, index *model.ForwardIndex) error {
	// 为特定索引获取加载锁
	loadMutex := fm.loadMutexForIndex(name)
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
func (fm *ForwardIndexManager) GetIndexByName(name string) (*model.ForwardIndex, bool) {
	fm.mutex.RLock()
	defer fm.mutex.RUnlock()

	index, exists := fm.indexes[name]
	return index, exists
}
