package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func HashUsingMD5(value string) string {
	hash := md5.Sum([]byte(value))
	hashedPass := hex.EncodeToString(hash[:])
	return hashedPass
}
