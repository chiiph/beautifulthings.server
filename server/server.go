package server

import (
	"beautifulthings/store"
	"time"

	"github.com/patrickmn/go-cache"
)

// TODO: implement the following
// SetPref(token, key, val) error or null if ok
// GetPrefs(token) error or json if ok

type Server interface {
	SignUp(b []byte) error
	SignIn(b []byte) ([]byte, error)
	Set(token string, date string, ct []byte) error
	Enumerate(token string, from, to string) ([]store.BeautifulThing, error)
}

type server struct {
	store   store.ObjectStore
	session *cache.Cache
}

func New(store store.ObjectStore) Server {
	return &server{
		store:   store,
		session: cache.New(24*time.Hour, 24*time.Hour),
	}
}
