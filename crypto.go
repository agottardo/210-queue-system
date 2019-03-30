package main

import "golang.org/x/crypto/bcrypt"

func GenerateSecretForCSid(csid string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(csid+config.AuthSecret), 14)
	return string(bytes)
}

func CheckSecretForCSid(secret string, csid string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(secret), []byte(csid+config.AuthSecret))
	return err == nil
}
