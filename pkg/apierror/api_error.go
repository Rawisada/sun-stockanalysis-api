package apierror

type ErrorCode string

const (
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeBadRequest     ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden      ErrorCode = "FORBIDDEN"
	ErrCodeConflict       ErrorCode = "CONFLICT"
	ErrCodeInternalError  ErrorCode = "INTERNAL_ERROR"
)

type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Status  int       `json:"status"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewNotFound(message string) *APIError {
	return &APIError{
		Code:    ErrCodeNotFound,
		Message: message,
		Status:  404,
	}
}

func NewBadRequest(message string) *APIError {
	return &APIError{
		Code:    ErrCodeBadRequest,
		Message: message,
		Status:  400,
	}
}

func NewUnauthorized(message string) *APIError {
	return &APIError{
		Code:    ErrCodeUnauthorized,
		Message: message,
		Status:  401,
	}
}

func NewForbidden(message string) *APIError {
	return &APIError{
		Code:    ErrCodeForbidden,
		Message: message,
		Status:  403,
	}
}

func NewConflict(message string) *APIError {
	return &APIError{
		Code:    ErrCodeConflict,
		Message: message,
		Status:  409,
	}
}

func NewInternalError(message string) *APIError {
	return &APIError{
		Code:    ErrCodeInternalError,
		Message: message,
		Status:  500,
	}
}
