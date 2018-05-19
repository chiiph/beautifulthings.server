package integration_tests

import (
	"beautifulthings/server"
	"beautifulthings/store"
	"testing"

	"beautifulthings/account"

	"github.com/stretchr/testify/require"
)

func set(t *testing.T, s server.Server, a *account.Account, token string, date, content string) {
	ct, err := a.Encrypt(content)
	require.NoError(t, err)
	err = s.Set(token, date, ct)
	require.NoError(t, err)
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
		m, err := a.Decrypt(ctit.Content)
		require.NoError(t, err)
		it := item{
			date:    ctit.Date.Format("2006-01-02"),
			content: m,
		}
		got = append(got, it)
	}
	return got
}

func testAddEnumerate(t *testing.T, s server.Server) {
	a := signup(t, s, "user1", "pass")

	token := signin(t, s, a)

	items := []item{
		{date: "2018-01-01", content: "item1"},
		{date: "2018-01-02", content: "item2"},
	}

	for _, it := range items {
		set(t, s, a, token, it.date, it.content)
	}

	got := enumerate(t, token, s, a, "2018-01-01", "2018-01-30")
	require.Len(t, got, len(items))
	require.Equal(t, items, got)
}

func TestAddEnumerate(t *testing.T) {
	s := server.New(store.NewInMemoryServer())
	testAddEnumerate(t, s)
}

func testAddEnumerateSkipOutside(t *testing.T, s server.Server) {
	a := signup(t, s, "user1", "pass")

	token := signin(t, s, a)

	items := []item{
		{date: "2018-01-01", content: "item1"},
		{date: "2018-01-02", content: "item2"},
		{date: "2018-01-10", content: "item3"},
		{date: "2018-02-10", content: "item4"},
	}

	for _, it := range items {
		set(t, s, a, token, it.date, it.content)
	}

	got := enumerate(t, token, s, a, "2018-01-01", "2018-01-30")
	require.Len(t, got, 3)
	require.Equal(t, items[:3], got)
}

func TestAddEnumerateSkipOutside(t *testing.T) {
	s := server.New(store.NewInMemoryServer())
	testAddEnumerateSkipOutside(t, s)
}
