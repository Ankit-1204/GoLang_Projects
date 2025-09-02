package internal

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

func (s *store) Get(key string) (interface{}, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	val, ok := s.mp[key]
	return val, ok
}

func (s *store) Set(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.mp[key] = value
}

func (s *store) Delete(key string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.mp, key)
	return nil
}
