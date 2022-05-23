package secure

import "golang.org/x/crypto/bcrypt"

//Hash receive string and add hash
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

//VerifyPassword compare password and hash and returns if they are the same
func VerifyPassword(passwordWithHash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordWithHash), []byte(password))
}
