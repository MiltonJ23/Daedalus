package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptCost is fixed at 12 according to the exigence NFR-SEC-02	in the SRS
const BcryptCost = 12

// DummyHash is a valied hash generated at start to guarantee a constant execution time
var DummyHash []byte

func init() {
	var err error
	DummyHash, err = bcrypt.GenerateFromPassword([]byte("dummy_timing_protection_string"), BcryptCost)
	if err != nil {
		panic("failed to generate dummy hash for timing attack protection: " + err.Error())
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func DummyCheck(password string) {
	_ = bcrypt.CompareHashAndPassword(DummyHash, []byte(password))
}
