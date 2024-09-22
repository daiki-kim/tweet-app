package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func Uint2String(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}

func String2Uint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 64)
	return uint(i)
}
