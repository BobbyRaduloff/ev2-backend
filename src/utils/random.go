package utils

import (
	"fmt"
	"math/rand"
)

func GetRandomString(length int) string {
	const characters = "abcdefghijklmnopqrstuvwxyz"
	var result string
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(characters))
		result += string(characters[randomIndex])
	}
	return result
}

func GetRandomEmail() string {
	usernameLength := 8

	username := GetRandomString(usernameLength)
	domain := "gmail.com"

	return fmt.Sprintf("%s@%s", username, domain)
}

func GetRandomEmailFromDomain(domain string) string {
	usernameLength := 8

	username := GetRandomString(usernameLength)

	return fmt.Sprintf("%s@%s", username, domain)
}
