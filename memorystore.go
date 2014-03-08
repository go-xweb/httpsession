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
	store.lock.Lock()
	node, ok := store.nodes[id]
	if !ok {
		node = &sessionNode{kvs: make(map[string]interface{}), expire: DefaultExpireTime}
		node.kvs[key] = value
		store.nodes[id] = node
		store.lock.Unlock()
	} else {
		store.lock.Unlock()
		node.lock.Lock()
		node.expire = DefaultExpireTime
		node.kvs[key] = value
		node.lock.Unlock()
	}
}

func (store *MemoryStore) Add(id Id) {
	node := &sessionNode{kvs: make(map[string]interface{}), expire: DefaultExpireTime}
	store.lock.Lock()
	store.nodes[id] = node
	store.lock.Unlock()
}

func (store *MemoryStore) Del(id Id, key string) bool {
	store.lock.RLock()
	node, ok := store.nodes[id]
	store.lock.RUnlock()
	if ok {
		node.lock.Lock()
		delete(node.kvs, key)
		node.lock.Unlock()
	}
	return true
}

func (store *MemoryStore) Exist(id Id) bool {
	store.lock.RLock()
	defer store.lock.RUnlock()
	_, ok := store.nodes[id]
	return ok
}

func (store *MemoryStore) Clear(id Id) bool {
	store.lock.RLock()
	defer store.lock.RUnlock()
	delete(store.nodes, id)
	return true
}

// TODO: gc
func (store *MemoryStore) Run() error {
	return nil
}
