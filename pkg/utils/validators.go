package utils

import (
	"errors"
	"regexp"
	"time"
)

var emailRx = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)

func ValidateEmail(email string) bool {

	return emailRx.MatchString(email)
}

func CalculateAge(dob time.Time) int {

	now := time.Now()
	age := now.Year() - dob.Year()

	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {

		age--
	}

	return age
}

func ValidateDOB(dob time.Time) error {

	if dob.After(time.Now()) {

		return errors.New("date_of_birth Cannot Be In The Future :( ")
	}

	return nil
}
