package controllers

import "net/http"

func NewNotFound(message string) *APIError {
	return &APIError{
		Code:    ErrCodeNotFound,
		Message: message,
		Status:  http.StatusNotFound,
	}
}

func NewBadRequest(message string) *APIError {
	return &APIError{
		Code:    ErrCodeBadRequest,
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

func NewUnauthorized(message string) *APIError {
	return &APIError{
		Code:    ErrCodeUnauthorized,
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

func NewForbidden(message string) *APIError {
	return &APIError{
		Code:    ErrCodeForbidden,
		Message: message,
		Status:  http.StatusForbidden,
	}
}

func NewInternalError(message string) *APIError {
	return &APIError{
		Code:    ErrCodeInternalError,
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}
