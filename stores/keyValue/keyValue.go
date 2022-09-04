package keyValue

import (
	"sync"
)

type KVStore_str struct {
	Store  map[string]string
	RwLock sync.RWMutex
}

type KVStore_num struct {
	Store  map[string]int32
	RwLock sync.RWMutex
}
