package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error, omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func Ok() Response {
	return Response{
		Status: StatusOk,
	}
}

func Error(message string) Response {
	return Response{
		Status: StatusError,
		Error:  message,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMessages []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMessages = append(errMessages, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMessages = append(errMessages, fmt.Sprintf("field %s is not a valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMessages, ", "),
	}
}
