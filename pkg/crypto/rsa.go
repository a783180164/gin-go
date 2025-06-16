package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
)

// GenerateRSAKeys generates a new RSA private key of specified bit size.
func GenerateRSAKeys(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

// EncryptWithPublicKey encrypts data with an RSA public key.
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) (string, error) {
	cipherData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, msg, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherData), nil
}

// DecryptWithPrivateKey decrypts base64 ciphertext with an RSA private key.
func DecryptWithPrivateKey(encoded string, priv *rsa.PrivateKey) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, data, nil)
}
