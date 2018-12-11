package libs

import (
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func TestGenToken(t *testing.T) {
	location, _ := time.LoadLocation("UTC")
	jsonWebToken := JsonWebToken{"secret-key",
		"", // Token
		jwt.StandardClaims{Subject: "1",
			ExpiresAt: time.Date(1985, 11, 8, 0, 0, 0, 0, location).Unix(),
		},
	}

	jsonWebToken.GenToken()
	if want := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjUwMDI1NjAwMCwic3ViIjoiMSJ9.ZvLQXxokez7zm_bdnoaGFX_mR-holMOthTQqb55-fh4"; jsonWebToken.Token != want {
		t.Errorf("jsonWebToken.Token == %s, want %s", jsonWebToken.Token, want)
	}
}
