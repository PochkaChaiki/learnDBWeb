package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const cost = 11

func SaltAndHashString(input string, salt string) (string, error) {
	output, err := bcrypt.GenerateFromPassword([]byte(input+salt), cost)
	return string(output), err
}
