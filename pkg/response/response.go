package response

import status "sun-stockanalysis-api/pkg/status"

type Status struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Remark  string `json:"remark"`
}

type ApiResponse[T any] struct {
	Status Status `json:"status"`
	Data   T      `json:"data"`
}

func Success[T any](data T) ApiResponse[T] {
	return ApiResponse[T]{
		Status: Status{
			Code:    status.CodeSuccess,
			Message: status.MsgSuccess,
			Remark:  "",
		},
		Data: data,
	}
}

func Error(code, message, remark string) ApiResponse[any] {
	return ApiResponse[any]{
		Status: Status{
			Code:    code,
			Message: message,
			Remark:  remark,
		},
		Data: nil,
	}
}
