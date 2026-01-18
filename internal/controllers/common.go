package controllers

import "strconv"

type Status struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Remark  *string `json:"remark"`
}

type DataResponse[T any] struct {
	Status Status `json:"status"`
	Data   T      `json:"data"`
}

type ErrorResponse struct {
	statusCode int `json:"-"`
	DataResponse[any]
}

func (e *ErrorResponse) Error() string {
	return e.Status.Message
}

func (e *ErrorResponse) GetStatus() int {
	return e.statusCode
}

func NewStatus(code string, message string, remark *string) Status {
	return Status{
		Code:    code,
		Message: message,
		Remark:  remark,
	}
}

func SuccessStatus() Status {
	return Status{
		Code:    "0",
		Message: "Success.",
		Remark:  nil,
	}
}

func InvalidStatus(message string) Status {
	if message == "" {
		message = "Invalid."
	}
	return Status{
		Code:    "400",
		Message: message,
		Remark:  nil,
	}
}

func NewDataResponse[T any](status Status, data T) DataResponse[T] {
	return DataResponse[T]{
		Status: status,
		Data:   data,
	}
}

func NewErrorResponse(status int, message string, errs ...error) *ErrorResponse {
	_ = errs
	return &ErrorResponse{
		statusCode: status,
		DataResponse: NewDataResponse[any](
			NewStatus(strconv.Itoa(status), message, nil),
			nil,
		),
	}
}
