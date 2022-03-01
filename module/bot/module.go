package bot

import (
	"botTele/module/bot/controller"
	"botTele/module/bot/service"
	"github.com/labstack/echo"
	"time"
)

var mBotController *controller.BotController

func Initialize(e *echo.Echo, timeout time.Duration) {
	botService := service.NewBotService(timeout)
	mBotController = controller.NewBotController(botService)

	go botService.AutoReply()

	initRouter(e)
}

func initRouter(e *echo.Echo) {
	e.POST("bot/send-chat", mBotController.SendChat)
}
