package integration_tests

import (
	"beautifulthings/account"
	"beautifulthings/server"
	"beautifulthings/store"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	os.Setenv("BT_INSECURE_KDF", "1")
}

func accBytes(t *testing.T, a *account.Account) []byte {
	b, err := a.Bytes()
	require.NoError(t, err)
	return b
}

func acc(t *testing.T, u, p string) *account.Account {
	a, err := account.New(u, p)
	require.NoError(t, err)
	return a
}

func signup(t *testing.T, s server.Server, u, p string) *account.Account {
	a := acc(t, u, p)
	b := accBytes(t, a)
	require.NoError(t, s.SignUp(b))
	return a
}

func TestBasicSignUp(t *testing.T) {
	s := server.New(store.NewInMemoryServer())
	signup(t, s, "user1", "pass")
}

func TestDoubleSignUpFails(t *testing.T) {
	s := server.New(store.NewInMemoryServer())
	a := signup(t, s, "user1", "pass")
	b := accBytes(t, a)
	require.Error(t, s.SignUp(b))
}

func TestMultiUserSignUp(t *testing.T) {
	s := server.New(store.NewInMemoryServer())
	signup(t, s, "user1", "pass")
	signup(t, s, "user2", "pass")
}

func TestBasicSignIn(t *testing.T) {
	s := server.New(store.NewInMemoryServer())
	a := signup(t, s, "user1", "pass")
	b := accBytes(t, a)

	cipherToken, err := s.SignIn(b)
	require.NoError(t, err)
	require.NotEmpty(t, cipherToken)

	token, err := a.Decrypt(cipherToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	ct, err := a.Encrypt("test1")
	require.NoError(t, err)
	err = s.Set(token, "2018-01-01", ct)
	require.NoError(t, err)

	err = s.Set(token+"42", "2018-01-01", ct)
	require.Error(t, err)
}
