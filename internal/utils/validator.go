package utils

import "regexp"

func ValidateEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

func ValidatePassword(password string) bool {
	return len(password) >= 8
}

func ValidateRequired(value string) bool {
	return len(value) > 0
}
