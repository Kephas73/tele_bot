package base_controller

import (
	"botTele/infrastructure/logger"
	"botTele/infrastructure/response_base"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

type BaseController struct{}

func (controller *BaseController) WriteSuccess(c echo.Context, v interface{}) error {
	res := response_base.Response{
		Message: "Success",
		Status:  http.StatusOK,
		Data:    v,
	}

	logger.Info("[API] - Return success: ", c.Request().RequestURI, time.Now())

	return c.JSON(http.StatusOK, res)
}

func (controller *BaseController) WriteSuccessEmptyContent(c echo.Context) error {
	res := response_base.Response{
		Message: "Success",
		Status:  http.StatusOK,
		Data:    nil,
	}

	logger.Info("[API] - Return success empty: ", res, c.Request().RequestURI, time.Now())

	return c.JSON(http.StatusOK, res)
}

func (controller *BaseController) writeError(c echo.Context, statusCode int, err response_base.ErrorResponse) error {
	res := response_base.Response{
		Message: "Failed",
		Status:  uint32(statusCode),
		Data:    err,
	}

	logger.Error("[API] - Response error: ", res, c.Request().RequestURI, time.Now())

	return c.JSON(statusCode, res)
}

func (controller *BaseController) WriteBadRequest(c echo.Context, errorRes response_base.ErrorResponse) error {
	return controller.writeError(c, http.StatusBadRequest, errorRes)
}
