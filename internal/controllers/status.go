package controllers

const (
	CodeSuccess        = "0"
	CodeInvalidParam   = "1001"
	CodeDataNotFound   = "1002"
	CodeUnauthorized   = "4001"
	CodeRequestTimeout = "4008"
	CodeSystemError    = "5001"
	CodeGeneralError   = "5002"

	MsgSuccess        = "Success."
	MsgInvalidParam   = "Invalid Param."
	MsgDataNotFound   = "Data Not Found."
	MsgUnauthorized   = "Unauthorized."
	MsgRequestTimeout = "Request Timeout"
	MsgSystemError    = "System Error."
	MsgGeneralError   = "General Error."
)

type ErrorCode string

const (
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeBadRequest     ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden      ErrorCode = "FORBIDDEN"
	ErrCodeConflict       ErrorCode = "CONFLICT"
	ErrCodeInternalError  ErrorCode = "INTERNAL_ERROR"
)