package pkg

import (
	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) []byte {
	salt := make([]byte, 16)
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return hash
}

func IsPasswordValid(hash []byte, password string) bool {
	passwordHash := HashPassword(password)

	if len(hash) != len(passwordHash) {
		return false
	}

	for i := 0; i < len(hash); i++ {
		if hash[i] != passwordHash[i] {
			return false
		}
	}

	return true
}
