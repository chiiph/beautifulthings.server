package store

import (
	"beautifulthings/account"

	"time"

	"github.com/patrickmn/go-cache"
)

func NewInMemory(a *account.Account) Store {
	return nil
}

func NewInMemoryServer() ServerStore {
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

func (m *memServerStore) GetForMonth(year int, month int, limit int) ([]byte, error) {
	panic("not implemented")
}

func (m *memServerStore) AddOrUpdate(date time.Time, b []byte) error {
	panic("not implemented")
}

func (m *memServerStore) Remove(date time.Time) error {
	panic("not implemented")
}

func (m *memServerStore) SetNagInterval(i string) error {
	panic("not implemented")
}
