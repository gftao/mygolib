package security

import (
	"crypto/hmac"
	"crypto/md5"
	"strings"
	"encoding/hex"
)

func HmacMd5(src, key string) string {
	k := []byte(key)
	mac := hmac.New(md5.New, k)
	mac.Write([]byte(src))
	cipherStr := mac.Sum(nil)
	return strings.ToUpper(hex.EncodeToString(cipherStr))
}

func VerifyHmacMd5(src, hmaced, key []byte) bool {
	mac := hmac.New(md5.New, key)
	mac.Write(src)
	expectedMAC := mac.Sum(nil)
	eMAC := strings.ToUpper(hex.EncodeToString(expectedMAC))
	re := strings.ToUpper(string(hmaced))
	return hmac.Equal([]byte(re), []byte(eMAC))
}
