package utils

import (
	"errors"
	"unicode"
)

func ValidatePassword(password string) error {
	var (
		hasMinLen	= false
		hasUpper	= false
		hasLower 	= false
		hasNumber	= false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasMinLen || !hasUpper || !hasLower || !hasNumber {
		return errors.New("password must be at least 8 characters long and contain at last one uppercase letter and one number")
	}

	return nil
}