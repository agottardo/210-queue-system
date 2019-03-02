package main

import "unicode"

func IsValidCSid(id string) bool {
	if len(id) != 4 && len(id) != 5 {
		return false
	}
	for i, char := range id {
		if i%2 == 0 && !unicode.IsLetter(char) {
			return false
		} else if i%2 != 0 && !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}