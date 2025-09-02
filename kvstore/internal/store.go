package store

import (
	"sync"
)

type store struct {
	lock sync.RWMutex
	mp   map[string]interface{}
}

func CreateStore() *store {
	return &store{mp: make(map[string]interface{})}
}

func (s *store) Get(key string) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.mp[key]
}

func (s *store) Set(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.mp[key] = value
}
