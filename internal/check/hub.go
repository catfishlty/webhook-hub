package check

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"regexp"
	"strconv"
)

func validateHttpMethod(method string) validation.MatchRule {
	return validation.Match(regexp.MustCompile(`$` + method + `^`)).
		ErrorObject(validation.NewError("validation_http_method", "http request method must be "+method))
}

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

func ValidateHttpRequest(method string, c *gin.Context) error {
	return validation.Validate(c.Request.Method, validation.Required, validation.When(method != "", validateHttpMethod(method)))
}
