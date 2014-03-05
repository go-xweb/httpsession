package httpsession

import (
	"sync"
	"time"
)

var _ Store = NewMemoryStore(30)

type sessionNode struct {
	lock   sync.RWMutex
	kvs    map[string]interface{}
	expire time.Duration
}

type MemoryStore struct {
	lock   sync.RWMutex
	nodes  map[Id]*sessionNode
	expire time.Duration
}

func NewMemoryStore(expire time.Duration) *MemoryStore {
	return &MemoryStore{nodes: make(map[Id]*sessionNode), expire: expire}
}

func (store *MemoryStore) Get(id Id, key string) interface{} {
	store.lock.RLock()
	node, ok := store.nodes[id]
	store.lock.RUnlock()
	if !ok {
		return nil
	}

	node.lock.Lock()
	node.expire = DefaultExpireTime
	node.lock.Unlock()

	node.lock.RLock()
	v, ok := node.kvs[key]
	node.lock.RUnlock()

	if !ok {
		return nil
	}
	return v
}

func (store *MemoryStore) Set(id Id, key string, value interface{}) {
	store.lock.RLock()
	node, ok := store.nodes[id]
	store.lock.RUnlock()
	if !ok {
		node = &sessionNode{kvs: make(map[string]interface{}), expire: DefaultExpireTime}
		node.kvs[key] = value
		store.lock.Lock()
		store.nodes[id] = node
		store.lock.Unlock()

		node.lock.Lock()

		node.lock.Unlock()
	} else {
		node.lock.Lock()
		node.expire = DefaultExpireTime
		node.kvs[key] = value
		node.lock.Unlock()
	}

}

func (store *MemoryStore) Del(id Id, key string) bool {
	return true
}

func (store *MemoryStore) DelAll(id Id) bool {
	return true
}

func (store *MemoryStore) Run() error {
	return nil
}
