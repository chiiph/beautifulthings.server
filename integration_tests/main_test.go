package integration_tests

import (
	"beautifulthings/server"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	os.Setenv("BT_INSECURE_KDF", "1")
}

func testBasicSignUp(t *testing.T, s server.Server) {
	signup(t, s, "user1", "pass")
}

func TestBasicSignUp(t *testing.T) {
	Run(t, testBasicSignUp)
}

func testDoubleSignUpFails(t *testing.T, s server.Server) {
	a := signup(t, s, "user1", "pass")
	b := accBytes(t, a)
	require.Error(t, s.SignUp(b))
}

func TestDoubleSignUpFails(t *testing.T) {
	Run(t, testDoubleSignUpFails)
}

func testMultiUserSignUp(t *testing.T, s server.Server) {
	signup(t, s, "user1", "pass")
	signup(t, s, "user2", "pass")
}

func TestMultiUserSignUp(t *testing.T) {
	Run(t, testMultiUserSignUp)
}

func testBasicSignIn(t *testing.T, s server.Server) {
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
func TestBasicSignIn(t *testing.T) {
	Run(t, testBasicSignIn)
}
