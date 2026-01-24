package controllers

type Status struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Remark  *string `json:"remark"`
}

type DataResponse[T any] struct {
	Status Status `json:"status"`
	Data   T      `json:"data"`
}

type APIError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Status  int         `json:"status"`
}

func (e *APIError) Error() string {
	return e.Message
}

type ErrorResponse struct {
	statusCode int `json:"-"`
	DataResponse[any]
}


type StatusResponse struct {
	Status int `status:"default"`
	Body   DataResponse[any]
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

func NewDataResponse[T any](status Status, data T) DataResponse[T] {
	return DataResponse[T]{
		Status: status,
		Data:   data,
	}
}

func NewGenericResponse[T any](status Status, data T) DataResponse[T] {
	return DataResponse[T]{
		Status: status,
		Data:   data,
	}
}

func SuccessResponse[T any](data T) DataResponse[T] {
	return NewDataResponse(NewStatus(CodeSuccess, MsgSuccess, nil), data)
}

func NewErrorResponse(httpStatus int, message string, errs ...error) *ErrorResponse {
	_ = errs
	status := statusForHTTP(httpStatus)
	if message != "" {
		status.Remark = &message
	}
	return &ErrorResponse{
		statusCode: httpStatus,
		DataResponse: NewDataResponse[any](status, nil),
	}
}
