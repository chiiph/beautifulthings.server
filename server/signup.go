package server

import (
	"beautifulthings/account"
	"beautifulthings/store"
	"crypto/rand"
	"encoding/base64"

	"github.com/pkg/errors"
)

var ErrAccountExists = errors.New("account already exists")
var ErrAccountDoesNotExist = errors.New("account does not exist")
var ErrInvalidSession = errors.New("invalid session")

func (s *server) accountExists(b []byte) (*account.Account, error) {
	a, err := account.FromBytes(b)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = a.Validate()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	b, err = s.store.Get(a.StorePath())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	a, err = account.FromBytes(b)

	// TODO: check that the pub key and the rest are the same as the stored version

	return a, err
}

func (s *server) SignUp(b []byte) error {
	_, err := s.accountExists(b)
	if err == nil {
		return ErrAccountExists
	}

	a, err := account.FromBytes(b)
	if err != nil {
		return ErrAccountExists
	}

	if err != nil && err != store.ErrNotFound {
		return errors.WithStack(err)
	}

	err = s.store.Set(a.StorePath(), b)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *server) token() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (s *server) setSession(token string, a *account.Account) {
	s.session.SetDefault(token, a)
}

func (s *server) validateSession(token string) (*account.Account, error) {
	rawAccount, found := s.session.Get(token)
	if !found {
		return nil, ErrInvalidSession
	}
	a := rawAccount.(*account.Account)
	return a, nil
}

func (s *server) SignIn(b []byte) ([]byte, error) {
	a, err := s.accountExists(b)
	if err != nil {
		return nil, ErrAccountDoesNotExist
	}

	token := s.token()
	c, err := a.Encrypt(token)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s.setSession(token, a)

	return c, nil
}
