package internal

import (
	"errors"
	"sync"
)

type Storage struct {
	mu sync.RWMutex
	db map[string][]byte
}

// NewStorage creates and returns a Storage instance.
func NewStorage() *Storage {
	m := make(map[string][]byte)
	return &Storage{
		db: m,
	}
}

// ReadToken servers as the storage layer GET function of /tokens endpoint,
// retrieving and returning the matching token from storage.
func (s *Storage) ReadToken(key string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.db[key]
	if !ok || val == nil {
		return nil, errors.New("unable to get token")
	}

	return val, nil
}

// WriteToken servers as the storage layer POST function of /tokens endpoint,
// creating a new token from storage.
func (s *Storage) WriteToken(token string, secret []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.db[token] = secret
	return nil
}

// UpdateToken servers as the storage layer PUT function of /tokens endpoint,
// changing the value of an existing token in storage.
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

// DeleteToken servers as the storage layer DELETE function of /tokens endpoint,
// removes the corresponding token from storage.
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