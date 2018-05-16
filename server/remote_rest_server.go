package server

import (
	"beautifulthings/store"
	"bytes"
	"net/http"

	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type remoteRestServer struct {
	addr string
}

func NewRemoteRest(addr string) Server {
	return &remoteRestServer{
		addr: addr,
	}
}

func (rs *remoteRestServer) SignUp(b []byte) error {
	resp, err := http.Post(rs.addr+"/signup", "application/octet-stream", bytes.NewReader(b))
	if err != nil {
		return errors.WithStack(err)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("Error signing up: %d", resp.StatusCode)
	}
	return nil
}

func (rs *remoteRestServer) SignIn(b []byte) ([]byte, error) {
	resp, err := http.Post(rs.addr+"/signin", "application/octet-stream", bytes.NewReader(b))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Error signing up: %d", resp.StatusCode)
	}
	sirb, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sir := &SignInResponse{}
	err = json.Unmarshal(sirb, sir)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return sir.EncryptedToken, nil
}

func (rs *remoteRestServer) Set(token string, date string, ct []byte) error {
	panic("not implemented")
}

func (rs *remoteRestServer) Enumerate(token string, from, to string) ([]store.BeautifulThing, error) {
	panic("not implemented")
}
