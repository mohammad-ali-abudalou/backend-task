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

	timeNow := time.Now()
	age := timeNow.Year() - dob.Year()

	if timeNow.Month() < dob.Month() || (timeNow.Month() == dob.Month() && timeNow.Day() < dob.Day()) {

		age--
	}

	return age
}

func ValidateDateOfBirth(dateOfBirth time.Time) error {

	if dateOfBirth.After(time.Now()) {

		return errors.New(ErrDateOfBirthCanNotInFuture)
	}

	return nil
}
