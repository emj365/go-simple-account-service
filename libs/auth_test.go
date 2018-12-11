package libs

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash := HashPassword("password", "salt")

	want := "b305cadbb3bce54f3aa59c64fec00dea"
	if hash != want {
		t.Errorf("hash == %q, want %q", hash, want)
	}
}
