package server

import (
	"beautifulthings/account"
	"beautifulthings/store"
	"beautifulthings/utils"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/nacl/box"
)

const dateLayout = "2006-01-02"
const contentLength = 240
const keyLength = 32

func (s *server) Set(token string, date string, ct []byte) error {
	a, err := s.validateSession(token)
	if err != nil {
		return errors.WithStack(err)
	}

	exactSize := len("2018-01-01")
	if len(date) != exactSize {
		return errors.Errorf("date too long: %d (max %d)", len(date), exactSize)
	}

	maxContent := contentLength + box.Overhead + keyLength
	if len(ct) > maxContent {
		return errors.Errorf("content too long: %d (max %d)", len(ct), maxContent)
	}

	_, err = time.Parse(dateLayout, date)
	if err != nil {
		return errors.WithStack(err)
	}

	err = s.store.Set(btPath(a, date), ct)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *server) Enumerate(token string, from, to string) ([]store.BeautifulThing, error) {
	a, err := s.validateSession(token)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fd, err := time.Parse(dateLayout, from)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	td, err := time.Parse(dateLayout, to)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !td.After(fd) {
		return nil, errors.New("from date needs to be before to date")
	}

	var list []store.BeautifulThing

	for d := fd; d.Before(td); d = d.AddDate(0, 0, 1) {
		ds := d.Format(dateLayout)
		path := btPath(a, ds)
		ct, err := s.store.Get(path)
		switch err {
		case store.ErrNotFound:
			continue
		case nil:
		default:
			return nil, errors.WithStack(err)
		}
		bt := store.BeautifulThing{
			Date:    d,
			Content: ct,
		}
		list = append(list, bt)
	}

	return list, nil
}

// TODO: maybe move to utils?
func btPath(a *account.Account, date string) string {
	return filepath.Join("beautifulthings", utils.S25664(a.Username), "1", date)
}
