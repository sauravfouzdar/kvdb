package storage

import "sync"

type Storage struct {
		data map[string]string
		mu   sync.RWMutex
}

func NewStorage() *Storage {
		return &Storage{
			data: make(map[string]string),
		}
}

/*
	RLock - multiple go routines can read(not write) at a time 
	Lock - only one go routine read/write at a time
*/
func (s *Storage) Get(key string) (string, bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()
		value, ok := s.data[key]
		return value, ok
}

func (s *Storage) Set(key, value string) {
		s.mu.Lock()
		defer s.mu.Unlock()

}

func (s *Storage) Delete(key string) {
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.data, key)
}