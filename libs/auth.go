package libs

import (
	"crypto/md5"
	"encoding/hex"
)

func HashPassword(rawPassword string, salt string) string {
	hasher := md5.New()
	hasher.Write([]byte(rawPassword + salt))
	password := hex.EncodeToString(hasher.Sum(nil))
	return password
}
