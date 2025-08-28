package utils

import (

    "regexp"
    "time"
    "errors"
)

func ValidateEmail(email string) bool {

    regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
    re := regexp.MustCompile(regex)
    return re.MatchString(email)
}

func CalculateAge(dob time.Time) int {

    now := time.Now()
    age := now.Year() - dob.Year()
    if now.YearDay() < dob.YearDay() {
	
        age--
    }
	
    return age
}

func ValidateDOB(dob time.Time) error {

    if dob.After(time.Now()) {
	
        return errors.New("date_of_birth cannot be in the future")
    }
	
    return nil
}
