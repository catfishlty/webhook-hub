package check

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"strconv"
)

func UsernamePassword(username, password string) error {
	return validation.Errors{
		"username": validation.Validate(username, validation.Required),
		"password": validation.Validate(password, validation.Required),
	}.Filter()
}

func ValidateId(id string) error {
	return validation.Validate(&id, validation.Required, is.UUIDv4)
}

func ValidatePage(page string) error {
	err := validation.Validate(&page, validation.Required, is.Int)
	if err != nil {
		return err
	}
	pageCount, _ := strconv.Atoi(page)
	return validation.Validate(pageCount, validation.Required, validation.Min(1))
}
