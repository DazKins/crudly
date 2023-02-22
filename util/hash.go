package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func StringHash(str string) string {
	bytes := sha256.Sum256([]byte(str))
	return hex.EncodeToString(bytes[:])
}
