package inmemory

import (
	"sync"
	"time"
)

type StorageTTL struct {
	mu      sync.RWMutex
	storage map[string]struct{}
	TTL     time.Duration
}

func NewStorageTTL(ttl time.Duration) *StorageTTL {
	return &StorageTTL{
		storage: make(map[string]struct{}),
		TTL:     ttl,
	}
}

func (s *StorageTTL) Add(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storage[key] = struct{}{}
	go func() {
		<-time.After(s.TTL)
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.storage, key)
	}()
}

func (s *StorageTTL) IsExist(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.storage[key]
	return ok
}
