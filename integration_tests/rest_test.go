package integration_tests

import (
	"beautifulthings/server"
	"beautifulthings/store"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

const addr = "localhost:8080"

func startRestServer(t *testing.T) func() {
	cancel, err := server.ServeRest(context.Background(), addr, store.NewInMemoryServer())
	require.NoError(t, err)
	return cancel
}

func TestRestBasicSignUp(t *testing.T) {
	cancel := startRestServer(t)
	defer cancel()
	s := server.NewRemoteRest("http://" + addr)
	a := signup(t, s, "user1", "pass")
	b := accBytes(t, a)
	require.Error(t, s.SignUp(b))
}

func TestRestBasicDoubleSign(t *testing.T) {
	cancel := startRestServer(t)
	defer cancel()
	s := server.NewRemoteRest("http://" + addr)
	a := signup(t, s, "user1", "pass")
	b := accBytes(t, a)
	require.Error(t, s.SignUp(b))
}

func TestRestBasicSignIn(t *testing.T) {
	cancel := startRestServer(t)
	defer cancel()
	s := server.NewRemoteRest("http://" + addr)
	a := signup(t, s, "user1", "pass")
	b := accBytes(t, a)

	cipherToken, err := s.SignIn(b)
	require.NoError(t, err)
	require.NotEmpty(t, cipherToken)

	token, err := a.Decrypt(cipherToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestRestAddEnumerate(t *testing.T) {
	cancel := startRestServer(t)
	defer cancel()
	s := server.NewRemoteRest("http://" + addr)
	testAddEnumerate(t, s)
}

func TestRestAddEnumerateSkipOutside(t *testing.T) {
	cancel := startRestServer(t)
	defer cancel()
	s := server.NewRemoteRest("http://" + addr)
	testAddEnumerateSkipOutside(t, s)
}
