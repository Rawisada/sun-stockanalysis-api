package apierror

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"sun-stockanalysis-api/pkg/response"
)

// HumaError is the shared error envelope for Huma responses.
type HumaError struct {
	Status response.Status `json:"status"`
	Data   any             `json:"data"`
	status int
}

func (e *HumaError) Error() string {
	return e.Status.Message
}

func (e *HumaError) GetStatus() int {
	return e.status
}

// NewHumaError creates a Huma StatusError using the business error envelope.
func NewHumaError(status int, msg string, errs ...error) huma.StatusError {
	businessStatus := statusForHTTP(status)
	remark := msg
	if len(errs) > 0 && errs[0] != nil {
		remark = errs[0].Error()
	}

	return &HumaError{
		Status: response.Status{
			Code:    businessStatus.Code,
			Message: businessStatus.Message,
			Remark:  remark,
		},
		Data:   nil,
		status: status,
	}
}

// ToHumaError maps application errors to the shared Huma error envelope.
func ToHumaError(err error) error {
	if err == nil {
		return nil
	}

	if statusErr, ok := err.(huma.StatusError); ok {
		return statusErr
	}

	if apiErr, ok := err.(*APIError); ok {
		status := statusForAPIError(apiErr)
		if status.Code == "" {
			status = statusForHTTP(apiErr.Status)
		}
		return &HumaError{
			Status: response.Status{
				Code:    status.Code,
				Message: status.Message,
				Remark:  err.Error(),
			},
			Data:   nil,
			status: apiErr.Status,
		}
	}

	if fiberErr, ok := err.(*fiber.Error); ok {
		status := statusForHTTP(fiberErr.Code)
		return &HumaError{
			Status: response.Status{
				Code:    status.Code,
				Message: status.Message,
				Remark:  err.Error(),
			},
			Data:   nil,
			status: fiberErr.Code,
		}
	}

	return NewHumaError(http.StatusInternalServerError, "Internal server error", err)
}
