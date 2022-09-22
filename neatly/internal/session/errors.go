package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type httpResponseHandler func(http.ResponseWriter, *http.Request) error

var (
	NotFoundError = NewError("not found", "", 404)
)

type Error struct {
	Err         error  `json:"-"`
	Code        int    `json:"code,omitempty"`
	Description string `json:"description"`
	Trace       string `json:"trace,omitempty"`
}

func Middleware(h httpResponseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error *Error
		err := h(w, r)
		if err != nil {
			if errors.As(err, &error) {
				w.WriteHeader(error.Code)
				w.Write(error.Marshal())
			}
		}
	}
}

func NewError(description, trace string, code int) *Error {
	return &Error{
		Err:         fmt.Errorf(description),
		Code:        code,
		Description: description,
		Trace:       trace,
	}
}

func (e *Error) What() string {
	return e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bytes
}

func UnauthorizedError(message string) *Error {
	return NewError(message, "", http.StatusUnauthorized)
}

func BadRequestError(message string) *Error {
	return NewError(message, "", http.StatusBadRequest)
}

//func internalError(trace string) *Error {
//	return NewError("system error", trace, -1)
//}
//
//func APIError(message, trace string, code int) *Error {
//	return NewError(message, trace, code)
//}
