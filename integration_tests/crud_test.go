package integration_tests

import (
	"beautifulthings/server"
	"testing"

	"github.com/stretchr/testify/require"
)

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
	Run(t, testAddEnumerate)
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
	Run(t, testAddEnumerateSkipOutside)
}

func testAddTooBig(t *testing.T, s server.Server) {
	a := signup(t, s, "user1", "pass")

	token := signin(t, s, a)

	maxContent := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" +
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" +
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" +
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

	set(t, s, a, token, "2018-01-01", maxContent)
	setFails(t, s, a, token, "2018-01-01", maxContent+"A")
	setFails(t, s, a, token, "2018-01-01 00:12", "A")
	setFails(t, s, a, token, "2018-01-1", "A")
}

func TestAddTooBig(t *testing.T) {
	Run(t, testAddTooBig)
}
