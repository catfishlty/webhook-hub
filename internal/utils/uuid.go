package utils

import (
	uuid "github.com/satori/go.uuid"
)

func UUID() string {
	return uuid.Must(uuid.NewV4(), nil).String()
}
