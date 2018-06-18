package server

import (
	"beautifulthings/store"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

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
		return nil, errors.Errorf("Error signing in: %d", resp.StatusCode)
	}
	sirb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()
	sir := &SignInResponse{}
	err = json.Unmarshal(sirb, sir)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return sir.EncryptedToken, nil
}

func (rs *remoteRestServer) urlWithToken(token string, section ...string) string {
	return fmt.Sprintf("%s/%s?token=%s", rs.addr, path.Join(section...), url.QueryEscape(token))
}

func (rs *remoteRestServer) Set(token string, date string, ct []byte) error {
	sr := SetRequest{
		Date: date,
		Ct:   ct,
	}
	b, err := json.Marshal(sr)
	if err != nil {
		return errors.WithStack(err)
	}

	resp, err := http.Post(
		rs.urlWithToken(token, "things"),
		"application/octet-stream",
		bytes.NewReader(b),
	)
	if err != nil {
		return errors.WithStack(err)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("Error setting beautiful thing: %d", resp.StatusCode)
	}
	return nil
}

func (rs *remoteRestServer) Enumerate(token string, from, to string) ([]store.BeautifulThing, error) {
	resp, err := http.Get(rs.urlWithToken(token, "things", from, to))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Error setting beautiful thing: %d", resp.StatusCode)
	}

	thingsb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	var things []store.BeautifulThing
	err = json.Unmarshal(thingsb, &things)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return things, nil
}
