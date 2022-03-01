package controller

import (
	"botTele/infrastructure/base_controller"
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
	return controller.WriteSuccess(c, controller.Service.Status(c))
}
