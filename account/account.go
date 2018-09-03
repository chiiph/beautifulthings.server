package account

import (
	"beautifulthings/utils"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

type Account struct {
	Username     string
	Pk           *[32]byte
	Sk           *[32]byte `json:"-"`
	Key          *[32]byte `json:"-"`
	EncryptedKey []byte
	Tz           string
	Offset       int
}

func New(username, password string) (*Account, error) {
	pk, sk, err := deriveKeyPair(username, password)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var key [32]byte
	_, err = rand.Read(key[:])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tz, offset := time.Now().Zone()

	return &Account{
		Username: username,
		Pk:       pk,
		Sk:       sk,
		Key:      &key,
		Tz:       tz,
		Offset:   offset,
	}, nil
}

func FromBytes(b []byte) (*Account, error) {
	a := &Account{}
	err := json.Unmarshal(b, a)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if a.Sk != nil {
		var key [32]byte
		rawKey, err := a.Decrypt(a.EncryptedKey)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		copy(key[:], rawKey)
		a.Key = &key
	}

	return a, nil
}

func (a *Account) Bytes() ([]byte, error) {
	if a.Sk != nil {
		encKey, err := a.Encrypt(string(a.Key[:]))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		a.EncryptedKey = encKey
	}

	return json.Marshal(a)
}

func ephNonce(epk, rpk *[32]byte) ([24]byte, error) {
	var nonce [24]byte
	h, err := blake2b.New256(nil)
	if err != nil {
		return [24]byte{}, errors.WithStack(err)
	}
	h.Write(epk[:])
	h.Write(rpk[:])

	nonceSlice := h.Sum(nil)
	copy(nonce[:], nonceSlice)
	return nonce, nil
}

func (a *Account) SymEncrypt(s string) ([]byte, error) {
	if a.Key == nil {
		return nil, errors.New("trying to encrypt without bootstrap")
	}
	var nonce [24]byte
	_, err := rand.Read(nonce[:])
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ct := secretbox.Seal(nonce[:], []byte(s), &nonce, a.Key)
	return ct, nil
}

func (a *Account) SymDecrypt(c []byte) (string, error) {
	if a.Key == nil {
		return "", errors.New("trying to decrypt without bootstrap")
	}
	var nonce [24]byte
	copy(nonce[:], c)
	c = c[24:]
	m, ok := secretbox.Open(nil, c, &nonce, a.Key)
	if !ok {
		return "", errors.New("error decrypting payload")
	}
	return string(m), nil
}

func (a *Account) Encrypt(s string) ([]byte, error) {
	epk, esk, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	nonce, err := ephNonce(epk, a.Pk)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	b := box.Seal(epk[:], []byte(s), &nonce, a.Pk, esk)
	return b, nil
}

func (a *Account) Decrypt(c []byte) (string, error) {
	if a.Sk == nil {
		panic("Trying to decrypt without a secret key")
	}
	var epk [32]byte
	copy(epk[:], c[:32])
	nonce, err := ephNonce(&epk, a.Pk)
	if err != nil {
		return "", errors.WithStack(err)
	}

	b, ok := box.Open(nil, c[32:], &nonce, &epk, a.Sk)
	if !ok {
		return "", errors.WithStack(errors.New("box.Open failed"))
	}

	return string(b), nil
}

func (a *Account) StorePath() string {
	return path.Join("account", utils.S25664(a.Username))
}

func (a *Account) Validate() error {
	if len(a.Username) == 0 {
		return errors.New("Username too short")
	}
	if a.Pk == nil {
		return errors.New("public key too short")
	}
	// TODO: validate timezone + offset
	return nil
}

func deriveKeyPair(username, password string) (pk *[32]byte, sk *[32]byte, err error) {
	seed, err := kdf(username, password)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	pk, sk, err = box.GenerateKey(bytes.NewReader(seed))
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return pk, sk, nil
}

func kdf(username, password string) ([]byte, error) {
	u := utils.S256(username)
	p := utils.S256(password)
	if os.Getenv("BT_INSECURE_KDF") != "" {
		return scrypt.Key(p, u, 2, 1, 1, 32)
	}

	return scrypt.Key(p, u, 1<<20, 8, 1, 32)
}
