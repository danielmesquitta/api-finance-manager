package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"github.com/danielmesquitta/api-finance-manager/internal/config/env"
	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	e *env.Env
}

func NewHasher(e *env.Env) *Hasher {
	return &Hasher{
		e: e,
	}
}

// Hash hashes a password using a secret key as additional security
func (h *Hasher) Hash(plaintext string) (string, error) {
	if plaintext == "" {
		return "", errors.New("input cannot be empty")
	}

	hash := hmac.New(sha256.New, []byte(h.e.HashSecretKey))
	hash.Write([]byte(plaintext))
	keyed := hash.Sum(nil)

	keyedStr := base64.StdEncoding.EncodeToString(keyed)

	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(keyedStr),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// Compare verifies if an input matches a hash that was created with a secret key
func (h *Hasher) Compare(
	hashed string,
	compare string,
) (bool, error) {
	hash := hmac.New(sha256.New, []byte(h.e.HashSecretKey))
	hash.Write([]byte(compare))
	keyed := hash.Sum(nil)

	keyedStr := base64.StdEncoding.EncodeToString(keyed)

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(keyedStr))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
