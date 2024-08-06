package validation

import (
	"errors"
	"observe/schema"
	"unicode"
)

// 1. username and password must not be empty
// 2. username should be at least 5 characters long
// 3. password should be at least 8 characters
// 4. password should contain at least one uppercase letter
// 5. password should contain at least one lowercase letter
// 6. password should contain at least one number

func ValidateUserForRegistration(user schema.User) error {
	if user.Username == "" {
		return errors.New("username must not be empty")
	}
	if user.Password == "" {
		return errors.New("password must not be empty")
	}
	if len(user.Username) < 5 {
		return errors.New("username must be at least 5 characters long")
	}
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	var hasUpper, hasLower, hasNumber bool
	for _, c := range user.Password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		}
	}
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	return nil
}

func ValidateUserForLogin(user schema.User) error {
	if user.Username == "" {
		return errors.New("username must not be empty")
	}
	if user.Password == "" {
		return errors.New("password must not be empty")
	}
	return nil
}
