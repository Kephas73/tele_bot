package response_base

import (
	"botTele/constant"
	"botTele/infrastructure/error_base"
)

type Response struct {
	Message string      `json:"Message"`
	Status  uint32      `json:"Status"`
	Data    interface{} `json:"Data"`
}

type ErrorResponse struct {
	ErrorCode uint32 `json:"ErrorCode"`
	Message   string `json:"Message"`
	Exception string `json:"Exception"`
}

func NewErrorResponse(err *error_base.Error) ErrorResponse {
	return ErrorResponse{
		ErrorCode: err.Code,
		Message:   err.Line,
		Exception: constant.LogErrorPrefix + err.Error.Error(),
	}
}
