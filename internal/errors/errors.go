package apperrors

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound   = errors.New("resource not found")
	ErrBadRequest = errors.New("bad request")
	ErrInternal   = errors.New("internal server error")
)

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

func StatusCode(err error) int {
	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, ErrBadRequest) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
