package integration_tests

import (
	"beautifulthings/account"
	"beautifulthings/server"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func set(t *testing.T, s server.Server, a *account.Account, token string, date, content string) {
	ct, err := a.SymEncrypt(content)
	require.NoError(t, err)
	err = s.Set(token, date, ct)
	require.NoError(t, err)
}

func setFails(t *testing.T, s server.Server, a *account.Account, token string, date, content string) {
	ct, err := a.SymEncrypt(content)
	require.NoError(t, err)
	err = s.Set(token, date, ct)
	require.Error(t, err)
}

func signin(t *testing.T, s server.Server, a *account.Account) string {
	b := accBytes(t, a)

	cipherToken, err := s.SignIn(b)
	require.NoError(t, err)
	require.NotEmpty(t, cipherToken)

	token, err := a.Decrypt(cipherToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	return token
}

type item struct {
	date    string
	content string
}

func enumerate(t *testing.T, token string, s server.Server, a *account.Account, from, to string) []item {
	res, err := s.Enumerate(token, from, to)
	require.NoError(t, err)

	var got []item
	for _, ctit := range res {
		m, err := a.SymDecrypt(ctit.Content)
		require.NoError(t, err)
		it := item{
			date:    ctit.Date.Format("2006-01-02"),
			content: m,
		}
		got = append(got, it)
	}
	return got
}
