package str

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash hash a string
func Hash(str string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		panic("Error on hashing password")
	}
	return string(bytes)
}

// CheckHash check hashed string with original string
func CheckHash(hashedStr string, str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(str))
	return err == nil
}
