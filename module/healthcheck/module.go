package healthcheck

import (
	"botTele/module/healthcheck/controller"
	"botTele/module/healthcheck/service"
	"time"

	"github.com/labstack/echo"
)

var mHealthCheckController *controller.HealthCheckController

func Initialize(e *echo.Echo, timeout time.Duration) {
	healthCheckService := service.NewHealthCheckService(timeout)
	mHealthCheckController = controller.NewHealthCheckController(healthCheckService)

	initRouter(e)
}

func initRouter(e *echo.Echo) {
	e.GET("bot/status", mHealthCheckController.Status)
}
