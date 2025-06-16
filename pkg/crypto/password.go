package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

const pepper = "key123131231313"

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+pepper), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 校验密码
func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+pepper))
	return err == nil
}
