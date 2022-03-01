package service

import (
	"github.com/labstack/echo"
	"time"
)

type IHealthCheckService interface {
	Status(c echo.Context) error
}

type HealthCheckService struct {
	Timeout time.Duration
}

func NewHealthCheckService(timeout time.Duration) IHealthCheckService {
	return &HealthCheckService{
		Timeout: timeout,
	}
}

func (service *HealthCheckService) Status(c echo.Context) error {
	return nil
}
