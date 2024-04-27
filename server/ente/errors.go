package ente

import "net/http"

type ErrorCode string

const (
	BadRequest    ErrorCode = "BAD_REQUEST"
	CONFLICT      ErrorCode = "CONFLICT"
	InternalError ErrorCode = "INTERNAL_ERROR"
)

type ApiError struct {
	Code           ErrorCode `json:"code"`
	Message        string    `json:"message"`
	HttpStatusCode int       `json:"-"`
}

func NewInternalError(message string) *ApiError {
	return &ApiError{
		Code:           InternalError,
		Message:        message,
		HttpStatusCode: http.StatusInternalServerError}
}

func (e *ApiError) NewError(message string) *ApiError {
	return &ApiError{
		Code:           e.Code,
		Message:        message,
		HttpStatusCode: e.HttpStatusCode}
}
