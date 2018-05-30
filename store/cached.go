package store

import (
	"github.com/hashicorp/golang-lru"
)

func NewCached(s ObjectStore) ObjectStore {
	l, err := lru.New(20000)
	if err != nil {
		panic(err)
	}
	return &cachedStore{
		lru: l,
		s:   s,
	}
}

type cachedStore struct {
	lru *lru.Cache
	s   ObjectStore
}

func (c *cachedStore) Get(url string) ([]byte, error) {
	v, found := c.lru.Get(url)
	if !found {
		return c.s.Get(url)
	}
	return v.([]byte), nil
}

func (c *cachedStore) Set(url string, val []byte) error {
	c.lru.Add(url, val)
	return c.s.Set(url, val)
}
