package utils

import (
	"errors"
	"regexp"
	"time"
)

var emailRx = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)

// ValidateEmail Checks If The Email Has A Valid Format :
func ValidateEmail(email string) bool {

	return emailRx.MatchString(email)
}

// CalculateAge Returns The Age In Years Based On Date Of Birth :
func CalculateAge(dob time.Time) int {

	now := time.Now()
	age := now.Year() - dob.Year()

	// Adjust If Birthday Hasn't Occurred Yet This Year.
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {

		age--
	}

	return age
}

// ValidateDateOfBirth Ensures The Date Is Not In The Future :
func ValidateDateOfBirth(dob time.Time) error {

	if dob.After(time.Now()) {

		return errors.New(ErrDateOfBirthCannotBeFuture.Error())
	}

	return nil
}
