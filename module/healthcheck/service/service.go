package service

import (
    "botTele/infrastructure/bot"
    "botTele/infrastructure/logger"
    "github.com/labstack/echo"
    "time"
)

type IHealthCheckService interface {
    Status(c echo.Context) error
}

type HealthCheckService struct {
    Bot     *bot.TelegramBot
    Timeout time.Duration
}

func NewHealthCheckService(timeout time.Duration) IHealthCheckService {
    return &HealthCheckService{
        Bot:     nil,
        Timeout: timeout,
    }
}

func (service *HealthCheckService) Status(c echo.Context) error {
    err := service.Bot.Status()
    if err != nil {
        logger.Error("HealthCheckService::SendChat: - Send chat error: %v", err)
    }
    return nil
}
