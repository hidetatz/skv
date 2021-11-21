package skv

import (
	"errors"
	"sync"
)

type Store struct {
	store map[string]string
	m     sync.Mutex
}

func NewStore() *Store {
	return &Store{store: map[string]string{}}
}

var ErrNotFound = errors.New("Not Found")

func (s *Store) Get(key string) (string, error) {
	// no need to lock
	v, ok := s.store[key]
	if !ok {
		return "", ErrNotFound
	}

	return v, nil
}

func (s *Store) Set(key, value string) {
	s.m.Lock()
	s.store[key] = value
	s.m.Unlock()
}
