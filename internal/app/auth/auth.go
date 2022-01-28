package auth

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

const secretKey = "my-super-secret-key"
const uidBytes = 16
const uidChars = uidBytes * 2

type contextKey string

var ContextUserIDKey = contextKey("uid")

func CheckUserIDCookieValue(value string) bool {
	if len(value) != 96 {
		return false
	}

	uid := value[:uidChars]

	sign, err := hex.DecodeString(value[uidChars:])
	if err != nil {
		return false
	}

	return hmac.Equal(signUserID(uid), sign)
}

func GetUserIDCookieValue(uid string) string {
	sign := signUserID(uid)

	return uid + hex.EncodeToString(sign)
}

// GenerateUserID returns 32-characters long hexadecimal string representing
// 16 random bytes (to be used as a unique user id)
func GenerateUserID() (string, error) {
	b, err := generateRandomBytes(uidBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func signUserID(uid string) []byte {
	secret := md5.Sum([]byte(secretKey))
	h := hmac.New(sha256.New, secret[:])
	h.Write([]byte(uid))
	sign := h.Sum(nil)

	return sign
}

func generateRandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}
