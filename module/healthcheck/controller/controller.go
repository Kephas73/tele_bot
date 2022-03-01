package controller

import (
    "botTele/infrastructure/base_controller"
    "botTele/infrastructure/error_base"
    "botTele/infrastructure/response_base"
    "botTele/module/healthcheck/service"
    "github.com/labstack/echo"
)

type HealthCheckController struct {
    base_controller.BaseController
    Service service.IHealthCheckService
}

func NewHealthCheckController(service service.IHealthCheckService) *HealthCheckController {
    return &HealthCheckController{
        Service: service,
    }
}

func (controller *HealthCheckController) Status(c echo.Context) error {
    err := controller.Service.Status(c)
    if err != nil {
        errApi := error_base.New(error_base.ErrorSendDataCode, err)
        resp := response_base.NewErrorResponse(errApi)
        return controller.WriteBadRequest(c, resp)
    }
    return controller.WriteSuccessEmptyContent(c)
}
