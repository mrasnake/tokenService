package internal

import (
	"errors"
	"sync"
)

type Storage struct {
	mu sync.RWMutex
	db map[string][]byte
}

func NewStorage() *Storage {
	m := make(map[string][]byte)
	return &Storage{
		db: m,
	}
}

// WriteToken creates a new object in storage.
func (s *Storage) WriteToken(token string, secret []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.db[token] = secret
	return nil
}

// ReadToken returns a matching object from storage.
func (s *Storage) ReadToken(key string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.db[key]
	if !ok || val == nil {
		return nil, errors.New("unable to get token")
	}

	return val, nil
}

// UpdateToken returns all objects in storage.
func (s *Storage) UpdateToken(token string, secret []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.db[token]
	if !ok || val == nil {
		return errors.New("token does not exist")
	}
	s.db[token] = secret
	return nil

}

// DeleteToken deletes the corresponding object from storage.
func (s *Storage) DeleteToken(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.db[key]
	if !ok || val == nil {
		return errors.New("cannot remove token that does not exist")
	}

	delete(s.db, key)
	val, ok = s.db[key]
	if ok || val != nil {
		return errors.New("unable to remove item")
	}
	return nil
}
