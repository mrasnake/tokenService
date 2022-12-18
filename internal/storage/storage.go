package datastore

import (
	"errors"
)

type Storage struct {
	db map[string][]byte
}

func NewStorage() *Storage {
	m := make(map[string][]byte)
	return &Storage{
		db: m,
	}
}

// WriteToken creates a new object in storage.
func (s *Storage) WriteToken(in string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.db[in] = true
	if !s.db[in] {
		return errors.New("item not added")
	}
	return nil
}

// ReadToken returns a matching object from storage.
func (s *Storage) ReadToken(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.db[key]
	if !ok || !val {
		return "", errors.New("unable to get item")
	}
	return key, nil
}

// DeleteToken deletes the corresponding object from storage.
func (s *Storage) DeleteToken(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.db[key]
	if !ok || !val {
		return errors.New("cannot remove item that does not exist")
	}

	delete(s.db, key)
	val, ok = s.db[key]
	if ok || val {
		return errors.New("unable to remove item")
	}
	return nil
}

// UpdateTokens returns all objects in storage.
func (s *Storage) UpdateTokens() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := []string{}
	for k := range s.db {
		keys = append(keys, k)
	}
	if len(keys) < 1 {
		return nil, errors.New("no items to list")
	}
	return keys, nil
}