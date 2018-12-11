package libs

import (
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func genToken() string {
	location, _ := time.LoadLocation("UTC")
	jsonWebToken := JsonWebToken{"secret-key",
		"", // Token
		jwt.StandardClaims{Subject: "1",
			ExpiresAt: time.Date(2099, 11, 8, 0, 0, 0, 0, location).Unix(),
		},
	}
	jsonWebToken.GenToken() // eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwOTc3NzkyMDAsInN1YiI6IjEifQ.IAeJrc-I5WDQaxpq1gKXvMoH2fHQf0APSY6U82jnl64
	return jsonWebToken.Token
}

func TestGenToken(t *testing.T) {
	token := genToken()
	if want := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQwOTc3NzkyMDAsInN1YiI6IjEifQ.IAeJrc-I5WDQaxpq1gKXvMoH2fHQf0APSY6U82jnl64"; token != want {
		t.Errorf("jsonWebToken.Token == %s, want %s", token, want)
	}
}

func TestDecode(t *testing.T) {
	token := genToken()
	jsonWebToken := JsonWebToken{"secret-key",
		token,
		jwt.StandardClaims{},
	}

	err := jsonWebToken.Decode()
	if want := error(nil); err != want {
		t.Errorf("jsonWebToken.Decode() == %s, want %v", err, want)
	}

	if want := "1"; jsonWebToken.Claims.Subject != want {
		t.Errorf("jsonWebToken.Token == %s, want %s", jsonWebToken.Claims.Subject, want)
	}

	if want := int64(4097779200); jsonWebToken.Claims.ExpiresAt != want {
		t.Errorf("jsonWebToken.Token == %v, want %v", jsonWebToken.Claims.ExpiresAt, want)
	}
}
