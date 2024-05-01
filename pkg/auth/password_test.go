package auth_test

import (
	"testing"

	"github.com/LambdaaTeam/Emenu/pkg/auth"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hash := auth.HashPassword(password)

	if len(hash) == 0 {
		t.Error("HashPassword() failed")
	}

	if string(hash) == password {
		t.Error("HashPassword() failed")
	}

	if len(hash) != len(auth.HashPassword(password)) {
		t.Error("HashPassword() failed")
	}
}
