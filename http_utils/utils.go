package http_utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gopybara/utils"
	"net/http"
	"reflect"
	"strings"
)

func ThrowHTTPError(ctx *gin.Context, err error, s ...int) {
	status := http.StatusInternalServerError
	errorMessage := err.Error()

	var e HttpError
	if errors.As(err, &e) {
		status = e.GetStatusCode()
		errorMessage = e.GetMessage()
	}

	if s != nil {
		status = s[0]
	}

	ctx.JSON(status, &ResponseError{
		Status: status,
		Error: HttpErrorObject{
			Message: errorMessage,
		},
	})
}

func GetValidationErrors(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "endwithts":
		return "URL should ends with trailing slash"
	case "oneof":
		return "Should be one of [" + strings.Join(strings.Split(fe.Param(), " "), ",") + "]"
	case "notempty":
		return "Param should not be empty"
	case "requiredIfEmpty":
		return "Param should not be empty if " + utils.FirstLetterToLower(fe.Param()) + " is empty"
	case "email":
		return "Param should be valid email " + fe.Param()
	case "url":
		return "Param should be valid url"
	case "min":
		return "Param should be greater than " + fe.Param()
	case "max":
		return "Param should be less than " + fe.Param()
	}
	return "Unknown error"
}

// ThrowValidationError is a function to throw validation error
// It accepts gin context, error and status code
// If status code is not provided, it will be 422
func ThrowValidationError(ctx *gin.Context, err error, s ...int) {
	if s == nil {
		s = []int{422}
	}
	status := s[0]
	var ve validator.ValidationErrors
	var out HttpErrorObject
	if errors.As(err, &ve) {
		out = HttpErrorObject{
			Name:    utils.Ptr("validation_error"),
			Message: err.Error(),
			Details: make([]*ErrorObjectDetails, 0),
		}
		for _, fe := range ve {
			out.Details = append(out.Details, &ErrorObjectDetails{
				Field: utils.FirstLetterToLower(fe.Field()),
				Issue: GetValidationErrors(fe),
			})
		}
	} else {
		out = HttpErrorObject{Message: err.Error()}
	}

	ctx.JSON(status, ResponseError{
		Status: status,
		Error:  out,
	})
}

type ResponseError struct {
	Status int             `json:"status"`
	Error  HttpErrorObject `json:"error"`
}

type HttpErrorObject struct {
	Name    *string               `json:"name,omitempty"`
	Message string                `json:"message"`
	Details []*ErrorObjectDetails `json:"details,omitempty"`
}

type ErrorObjectDetails struct {
	Issue string `json:"issue,omitempty"`
	Field string `json:"field,omitempty"`
}

type HTTPResponse struct {
	Status int         `json:"status" example:"200"`
	Data   interface{} `json:"data,omitempty"`
	Meta   interface{} `json:"meta,omitempty"`
}

func NewHttpResponse(status int, data interface{}, meta ...interface{}) (int, *HTTPResponse) {
	responseMeta := interface{}(nil)
	if len(meta) > 0 {
		responseMeta = meta[0]
	}
	response := &HTTPResponse{
		Status: status,
		Data:   data,
		Meta:   responseMeta,
	}

	elem := reflect.ValueOf(data)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}

	if meta == nil && (data != nil && elem.Kind() == reflect.Slice || elem.Kind() == reflect.Array) {
		response.Meta = MetaCollectionTotals{
			Total: reflect.ValueOf(data).Len(),
		}
	}

	return status, response
}

type MetaCollectionTotals struct {
	Total int `json:"total" example:"6"`
}
