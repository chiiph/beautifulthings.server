package store

import (
	"beautifulthings/account"

	"github.com/patrickmn/go-cache"
)

func NewInMemory(a *account.Account) ObjectStore {
	return nil
}

func NewInMemoryServer() ObjectStore {
	return &memServerStore{
		c: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

type memServerStore struct {
	c *cache.Cache
}

func (m *memServerStore) Get(url string) ([]byte, error) {
	v, found := m.c.Get(url)
	if !found {
		return nil, ErrNotFound
	}
	return v.([]byte), nil
}

func (m *memServerStore) Set(url string, val []byte) error {
	m.c.SetDefault(url, val)
	return nil
}
