package bot

import (
    "botTele/module/bot/controller"
    "botTele/module/bot/service"
    "github.com/labstack/echo"
    "time"
)

var mBotController *controller.BotController
var BotServiceGlobal service.IBotService

func Initialize(e *echo.Echo, timeout time.Duration) {
    botService := service.NewBotService(timeout)
    mBotController = controller.NewBotController(botService)
    BotServiceGlobal = botService

    /* go func() {
        botService.WorkerUploadFile()
        time.Sleep(time.Minute * 5)
    }()*/

    initRouter(e)
}

func initRouter(e *echo.Echo) {
    e.POST("bot/send-chat", mBotController.SendChat)

    e.POST("file/upload", mBotController.UploadFile)

    e.POST("ip/init", mBotController.InitIP)
    e.GET("ip/random", mBotController.RandomIP)
    e.GET("class", mBotController.GetClass)

    e.GET("queryI", mBotController.ListI)
}
