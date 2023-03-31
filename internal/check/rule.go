package check

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"regexp"
	"strings"
)

var httpMethods = []string{
	http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut,
	http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace,
}

var logLevels = []string{
	"panic", "fatal", "error", "warn", "warning", "info", "debug", "trace",
}

var dbTypes = []string{
	"mysql", "postgres", "sqlite", "sqlserver",
}

var (
	IsKey = validation.Match(regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")).
		ErrorObject(validation.NewError("validation_is_key", "must contain English letters and digits and special characters '!#$%&'*+/=?^_`{|}~.-' only"))
	IsHttpMethod = convertToInRule(httpMethods).
		ErrorObject(validation.NewError("validation_is_http_method", fmt.Sprintf("must be a valid http method(%s)", strings.Join(httpMethods, ", "))))
	IsLogLevel = convertToInRule(logLevels).
		ErrorObject(validation.NewError("validation_is_log_level", fmt.Sprintf("must be a valid log level(%s)", strings.Join(logLevels, ", "))))
	IsDbType = convertToInRule(dbTypes).
		ErrorObject(validation.NewError("validation_is_db_type", fmt.Sprintf("must be a valid database type(%s)", strings.Join(dbTypes, ", "))))
)

func convertToInRule[T any](l []T) validation.InRule {
	t := make([]any, len(l))
	for i, v := range l {
		t[i] = v
	}
	return validation.In(t...)
}
