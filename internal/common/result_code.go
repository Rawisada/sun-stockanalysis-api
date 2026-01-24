package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ToResponse(err error) (int, DataResponse[any]) {
	if err == nil {
		return http.StatusOK, NewGenericResponse[any](
			NewStatus(CodeSuccess, MsgSuccess, nil),
			nil,
		)
	}

	if apiErr, ok := err.(*APIError); ok {
		status := statusForAPIError(apiErr)
		if status.Code == "" {
			status = statusForHTTP(apiErr.Status)
		}
		remark := err.Error()
		status.Remark = &remark
		return apiErr.Status, NewDataResponse[any](status, nil)
	}

	if fiberErr, ok := err.(*fiber.Error); ok {
		status := statusForHTTP(fiberErr.Code)
		remark := err.Error()
		status.Remark = &remark
		return fiberErr.Code, NewDataResponse[any](status, nil)
	}

	status := statusForHTTP(http.StatusInternalServerError)
	remark := err.Error()
	status.Remark = &remark
	return http.StatusInternalServerError, NewDataResponse[any](status, nil)
}

func statusForAPIError(apiErr *APIError) Status {
	switch apiErr.Code {
	case ErrCodeNotFound:
		return Status{Code: CodeDataNotFound, Message: MsgDataNotFound}
	case ErrCodeBadRequest:
		return Status{Code: CodeInvalidParam, Message: MsgInvalidParam}
	case ErrCodeUnauthorized:
		return Status{Code: CodeUnauthorized, Message: MsgUnauthorized}
	case ErrCodeForbidden:
		return Status{Code: CodeUnauthorized, Message: MsgUnauthorized}
	case ErrCodeConflict:
		return Status{Code: CodeInvalidParam, Message: MsgInvalidParam}
	case ErrCodeInternalError:
		return Status{Code: CodeSystemError, Message: MsgSystemError}
	default:
		return Status{}
	}
}

func statusForHTTP(httpStatus int) Status {
	switch httpStatus {
	case http.StatusBadRequest:
		return Status{Code: CodeInvalidParam, Message: MsgInvalidParam}
	case http.StatusUnauthorized:
		return Status{Code: CodeUnauthorized, Message: MsgUnauthorized}
	case http.StatusForbidden:
		return Status{Code: CodeUnauthorized, Message: MsgUnauthorized}
	case http.StatusNotFound:
		return Status{Code: CodeDataNotFound, Message: MsgDataNotFound}
	case http.StatusConflict:
		return Status{Code: CodeInvalidParam, Message: MsgInvalidParam}
	case http.StatusRequestTimeout:
		return Status{Code: CodeRequestTimeout, Message: MsgRequestTimeout}
	case http.StatusInternalServerError:
		return Status{Code: CodeSystemError, Message: MsgSystemError}
	default:
		return Status{Code: CodeGeneralError, Message: MsgGeneralError}
	}
}
