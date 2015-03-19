package utils

import "regexp"

func IsValidEmail(email string) bool {
	m, err := regexp.MatchString("[a-zA-Z0-9\\.\\+_-]+@[a-zA-Z0-9\\.\\+_-]+\\.[a-zA-Z0-9]+", email)
	if err != nil {
		return false
	}
	return m
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if a == b {
			return true
		}
	}
	return false
}
