package utils

import (
	"fmt"
)

const (
	ErrAuthUser = "Error authorizing user: "
)

type ErrMessage struct {
	Code   int    `json:"code"`
	Message string `json:"message"`
}

func (e ErrMessage) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func CustomError() error {
	return ErrMessage{
		Code:    500,
		Message: "Error",
	}
}
