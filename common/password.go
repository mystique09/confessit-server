package common

import "golang.org/x/crypto/bcrypt"

func CheckPassword(encryptedPass, plainText []byte) error {
	return bcrypt.CompareHashAndPassword(encryptedPass, plainText)
}

func HashPassword(plainText string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)

	return string(hashedPass), err
}
