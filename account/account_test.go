package account

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	os.Setenv("BT_INSECURE_KDF", "1")
}

func TestAccount_Bytes(t *testing.T) {
	a1, e := New("user", "pass")
	require.NoError(t, e)
	require.NotEmpty(t, a1.Username)
	require.NotEmpty(t, a1.Sk)
	require.NotEmpty(t, a1.Pk)
	require.NotEmpty(t, a1.Offset)
	require.NotEmpty(t, a1.Tz)

	b, e := a1.Bytes()
	require.NoError(t, e)

	a2, e := FromBytes(b)
	require.NoError(t, e)
	// FromBytes is only used in the server, so we don't have a Sk
	a1.Sk = nil
	a1.Key = nil

	require.Equal(t, a1, a2)
}

func TestAccount_EncryptDecryptFromServer(t *testing.T) {
	a1, e := New("user", "pass")
	require.NoError(t, e)

	// As it would be gotten on the server
	b, e := a1.Bytes()
	require.NoError(t, e)
	a2, e := FromBytes(b)
	require.NoError(t, e)

	m := "some string"
	ct, e := a2.Encrypt(m)
	require.NoError(t, e)

	me, e := a1.Decrypt(ct)
	require.NoError(t, e)
	require.Equal(t, m, me)
}

func TestAccount_EncryptDecryptFromSameAccount(t *testing.T) {
	a1, e := New("user", "pass")
	require.NoError(t, e)

	m := "some string"
	ct, e := a1.Encrypt(m)
	require.NoError(t, e)

	me, e := a1.Decrypt(ct)
	require.NoError(t, e)
	require.Equal(t, m, me)
}

func TestAccount_EncryptDecryptFailsOnModify(t *testing.T) {
	a1, e := New("user", "pass")
	require.NoError(t, e)

	m := "some string"
	ct, e := a1.Encrypt(m)
	require.NoError(t, e)

	ct[2] = ct[2] + '1'

	_, e = a1.Decrypt(ct)
	require.Error(t, e)
}

func TestAccount_Validate(t *testing.T) {
	a, e := New("user", "pass")
	require.NoError(t, e)
	require.NoError(t, a.Validate())
	a.Sk = nil
	require.NoError(t, a.Validate())
	a.Pk = nil
	require.Error(t, a.Validate())

	a, e = New("", "")
	require.NoError(t, e)
	require.Error(t, a.Validate())
}

func TestAccount_Sym(t *testing.T) {
	a, e := New("user", "pass")
	require.NoError(t, e)
	require.NotZero(t, a.Key)

	s := "test"
	ct, err := a.SymEncrypt(s)
	require.NoError(t, err)
	require.NotZero(t, ct)

	m, err := a.SymDecrypt(ct)
	require.NoError(t, err)
	require.Equal(t, s, m)

	ct2, err := a.SymEncrypt(s)
	require.NoError(t, err)
	require.NotZero(t, ct2)
	require.NotEqual(t, ct, ct2)

	m2, err := a.SymDecrypt(ct2)
	require.NoError(t, err)
	require.Equal(t, s, m2)

	ct2 = append(ct2, byte(' '))
	_, err = a.SymDecrypt(ct2)
	require.Error(t, err)
}
