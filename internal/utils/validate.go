package utils

import (
	"net/mail"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/grozaqueen/julse/internal/errs"
)

func ValidateRegistration(email string, username string, password string, repeatPassword string) error {
	if err := ValidateEmailAndPassword(email, password); err != nil {
		return err
	}

	if password != repeatPassword {
		return errs.PasswordsDoNotMatch
	}

	if !IsValidUsername(username) {
		return errs.InvalidUsernameFormat
	}

	return nil
}

func ValidateEmailAndPassword(email string, password string) error {
	switch {
	case !IsValidEmail(email):
		return errs.InvalidEmailFormat
	case !isValidPassword(password):
		return errs.InvalidPasswordFormat
	}
	return nil
}

func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" || len(email) > 254 {
		return false
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	if addr.Address != email {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]
	if domain == "" {
		return false
	}

	if !strings.Contains(domain, ".") {
		return false
	}
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") || strings.Contains(domain, "..") {
		return false
	}

	return true
}

func IsValidUsername(username string) bool {
	username = strings.TrimSpace(username)

	n := utf8.RuneCountInString(username)
	if n < 2 || n > 40 {
		return false
	}

	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func isValidPassword(password string) bool {
	hasMinLen := utf8.RuneCountInString(password) >= 5
	var hasNumber, hasLower bool

	for _, r := range password {
		switch {
		case unicode.IsNumber(r):
			hasNumber = true
		case unicode.IsLower(r):
			hasLower = true
		}
	}
	return hasMinLen && hasNumber && hasLower
}
