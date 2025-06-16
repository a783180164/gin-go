package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// SignHMAC computes an HMAC-SHA256 signature of data.
func SignHMAC(data, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// VerifyHMAC checks that signature matches data.
func VerifyHMAC(data, key []byte, signature string) bool {
	expected := SignHMAC(data, key)
	return hmac.Equal([]byte(expected), []byte(signature))
}
