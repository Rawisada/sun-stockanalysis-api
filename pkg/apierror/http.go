package apierror

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	status "sun-stockanalysis-api/pkg/status"
	"sun-stockanalysis-api/pkg/response"
)

// ToResponse maps application errors to a consistent HTTP status and response envelope.
func ToResponse(err error) (int, response.ApiResponse[any]) {
	if err == nil {
		return http.StatusOK, response.Success[any](nil)
	}

	if apiErr, ok := err.(*APIError); ok {
		status := statusForAPIError(apiErr)
		if status.Code == "" {
			status = statusForHTTP(apiErr.Status)
		}
		status.Remark = err.Error()
		return apiErr.Status, response.Error(status.Code, status.Message, status.Remark)
	}

	if fiberErr, ok := err.(*fiber.Error); ok {
		status := statusForHTTP(fiberErr.Code)
		status.Remark = err.Error()
		return fiberErr.Code, response.Error(status.Code, status.Message, status.Remark)
	}

	status := statusForHTTP(http.StatusInternalServerError)
	status.Remark = err.Error()
	return http.StatusInternalServerError, response.Error(status.Code, status.Message, status.Remark)
}

func statusForAPIError(apiErr *APIError) response.Status {
	switch apiErr.Code {
	case ErrCodeNotFound:
		return response.Status{Code: status.CodeDataNotFound, Message: status.MsgDataNotFound}
	case ErrCodeBadRequest:
		return response.Status{Code: status.CodeInvalidParam, Message: status.MsgInvalidParam}
	case ErrCodeUnauthorized:
		return response.Status{Code: status.CodeUnauthorized, Message: status.MsgUnauthorized}
	case ErrCodeForbidden:
		return response.Status{Code: status.CodeUnauthorized, Message: status.MsgUnauthorized}
	case ErrCodeConflict:
		return response.Status{Code: status.CodeInvalidParam, Message: status.MsgInvalidParam}
	case ErrCodeInternalError:
		return response.Status{Code: status.CodeSystemError, Message: status.MsgSystemError}
	default:
		return response.Status{}
	}
}

func statusForHTTP(httpStatus int) response.Status {
	switch httpStatus {
	case http.StatusBadRequest:
		return response.Status{Code: status.CodeInvalidParam, Message: status.MsgInvalidParam}
	case http.StatusUnauthorized:
		return response.Status{Code: status.CodeUnauthorized, Message: status.MsgUnauthorized}
	case http.StatusForbidden:
		return response.Status{Code: status.CodeUnauthorized, Message: status.MsgUnauthorized}
	case http.StatusNotFound:
		return response.Status{Code: status.CodeDataNotFound, Message: status.MsgDataNotFound}
	case http.StatusConflict:
		return response.Status{Code: status.CodeInvalidParam, Message: status.MsgInvalidParam}
	case http.StatusRequestTimeout:
		return response.Status{Code: status.CodeRequestTimeout, Message: status.MsgRequestTimeout}
	case http.StatusInternalServerError:
		return response.Status{Code: status.CodeSystemError, Message: status.MsgSystemError}
	default:
		return response.Status{Code: status.CodeGeneralError, Message: status.MsgGeneralError}
	}
}
