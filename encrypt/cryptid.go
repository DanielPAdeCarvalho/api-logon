package encrypt

import (
	"login-app/utils"

	"golang.org/x/crypto/bcrypt"
)

func EncrytpHash(senha string, logs utils.Loggar) string {
	senhaB := []byte(senha)
	senhaH, err := bcrypt.GenerateFromPassword(senhaB, bcrypt.MinCost)
	utils.Check(err, logs)
	senhaS := string(senhaH)
	return senhaS
}

func CheckHash(senha string, hash string, logs utils.Loggar) bool {
	senhaB := []byte(senha)
	hashB := []byte(hash)
	err := bcrypt.CompareHashAndPassword(hashB, senhaB)
	if err == nil {
		return true
	}
	utils.Check(err, logs)
	return false
}
