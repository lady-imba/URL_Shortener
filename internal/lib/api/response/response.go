package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

type Response struct {
	Status string `json:"status"` //Error Ok
	Error string `json:"error,omitempty"`
}

const (
	StatusOK = "ОК"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error: msg,
	}
}

func ValidationError(errors validator.ValidationErrors) Response {
	var errMsgs []string 

	for _, err := range errors {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error: strings.Join(errMsgs, ", "),
	}
}

