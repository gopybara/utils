package http_utils

import (
	"errors"
	"net/http"
)

var (
	ErrCacheNotFound         = NewHTTPError(http.StatusNotFound, "cache not found")
	ErrorNotFound            = NewHTTPError(http.StatusNotFound, "not found")
	ErrorUnauthorized        = NewHTTPError(http.StatusUnauthorized, "unauthorized")
	ErrorUserNotFound        = NewHTTPError(http.StatusNotFound, "user not found")
	ErrorPaymentRequired     = NewHTTPError(http.StatusPaymentRequired, "payment required")
	ErrorImageRequired       = NewHTTPError(http.StatusBadRequest, "image required")
	ErrorInternalServerError = NewHTTPError(http.StatusInternalServerError, "internal server error")
	ErrorTooManyRequests     = NewHTTPError(http.StatusTooManyRequests, "too many requests")
)

type HttpError struct {
	error
	statusCode      int
	frontendMessage string
}

func (e HttpError) GetStatusCode() int {
	return e.statusCode
}

func (e HttpError) GetMessage() string {
	if e.frontendMessage != "" {
		return e.frontendMessage
	}

	return e.Error()
}

func NewHTTPError(statusCode int, msg string, frontendMessage ...string) HttpError {
	hr := HttpError{error: errors.New(msg), statusCode: statusCode}
	if len(frontendMessage) > 0 {
		hr.frontendMessage = frontendMessage[0]
	}

	return hr
}

func WrapHTTPError(statusCode int, err error, frontendMessage ...string) HttpError {
	hr := HttpError{error: err, statusCode: statusCode}
	if len(frontendMessage) > 0 {
		hr.frontendMessage = frontendMessage[0]
	}

	return hr
}
