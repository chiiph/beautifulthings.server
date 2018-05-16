package store

import (
	"beautifulthings/account"
	"time"

	"github.com/pkg/errors"
)

type BeautifulThing struct {
	Date    time.Time
	Content []byte
	Extras  map[string][]string
}

// Store is something that the client uses, it calls the server side APIs
type Store interface {
	GetForMonth(year int, month int, limit int) ([]byte, error)

	AddOrUpdate(date time.Time, b []byte) error
	Remove(date time.Time) error

	// TODO: Add a way to store extras

	SetNagInterval(i string) error
}

var ErrNotFound = errors.New("does not exist")

// ServerStore is an enhanced Store that is used by the server. The enhancement is used for account logic
type ServerStore interface {
	Store

	Get(url string) ([]byte, error)
	Set(url string, val []byte) error
}

func NewRemote(a *account.Account) Store {
	return nil
}
