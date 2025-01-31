package models

import "errors"

var (
	ErrNotFound = errors.New("not found")
)

type MemStore struct {
	list map[string]User
}

func NewMemStore() *MemStore {
	list := make(map[string]User)
	return &MemStore{
		list,
	}
}

func (m MemStore) Add(name string, user User) error {
	m.list[name] = user
	return nil
}

func (m MemStore) Get(name string) (User, error) {

	if val, ok := m.list[name]; ok {
		return val, nil
	}

	return User{}, ErrNotFound
}

func (m MemStore) List() (map[string]User, error) {
	return m.list, nil
}

func (m MemStore) Update(name string, user User) error {

	if _, ok := m.list[name]; ok {
		m.list[name] = user
		return nil
	}

	return ErrNotFound
}

func (m MemStore) Remove(name string) error {
	delete(m.list, name)
	return nil
}
