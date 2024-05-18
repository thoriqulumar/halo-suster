package controller

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"strconv"
	"time"
)

func validateNIP(fl validator.FieldLevel) bool {
	nip := fl.Field().Int()
	nipStr := strconv.FormatInt(nip, 10)

	// Check first three digits
	if nipStr[:3] != "615" && nipStr[:3] != "303" {
		return false
	}

	// Check fourth digit (gender)
	if nipStr[3] != '1' && nipStr[3] != '2' {
		return false
	}

	// Check year
	year, _ := strconv.Atoi(nipStr[4:8])
	currentYear := time.Now().Year()
	if year < 2000 || year > currentYear {
		return false
	}

	// Check month
	month, _ := strconv.Atoi(nipStr[8:10])
	if month < 1 || month > 12 {
		return false
	}

	return true
}

func customURL(fl validator.FieldLevel) bool {
	// Regular expression pattern for a valid URL
	// This pattern requires the URL to start with http:// or https://
	// followed by a valid domain name
	pattern := `^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(\/\S*)?$`

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Match the URL against the regular expression
	url := fl.Field().String()
	return regex.MatchString(url)
}
