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

type BootstrapPayload struct {
	EncryptedKey []byte
}

var ErrNotFound = errors.New("does not exist")

// ObjectStore is an enhanced Store that is used by the server. The enhancement is used for account logic
type ObjectStore interface {
	Get(url string) ([]byte, error)
	Set(url string, val []byte) error
}

func NewRemote(a *account.Account) ObjectStore {
	return nil
}
