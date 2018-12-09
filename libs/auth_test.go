package libs

import (
	"log"
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash := HashPassword("password", "salt")

	if hash != "b305cadbb3bce54f3aa59c64fec00dea" {
		log.Fatalf("TestGenPasswordHash ger wrong hash, which is %s", hash)
	}
}
