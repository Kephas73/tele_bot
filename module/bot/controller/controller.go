package controller

import (
    "botTele/constant"
    "botTele/infrastructure/base_controller"
    "botTele/infrastructure/error_base"
    "botTele/infrastructure/logger"
    "botTele/infrastructure/response_base"
    "botTele/model"
    "botTele/module/bot/service"
    "context"
    "fmt"
    "github.com/labstack/echo"
)

type BotController struct {
    base_controller.BaseController
    Service service.IBotService
}

func NewBotController(service service.IBotService) *BotController {
    return &BotController{
        Service: service,
    }
}

func (controller *BotController) SendChat(c echo.Context) error {

    var data model.RawData

    logger.Info("BotController:SendChat:  OK")

    err := c.Bind(&data)
    if err != nil {
        errApi := error_base.New(error_base.ErrorBindDataCode, err)
        resp := response_base.NewErrorResponse(errApi)
        return controller.WriteBadRequest(c, resp)
    }

    if len(data.Text) == constant.ValueEmpty && data.Object == nil {
        errApi := error_base.New(error_base.ErrorValidDataCode, fmt.Errorf("data empty"))
        resp := response_base.NewErrorResponse(errApi)
        return controller.WriteBadRequest(c, resp)
    }

    err = controller.Service.SendChat(data)
    if err != nil {
        errApi := error_base.New(error_base.ErrorSendDataCode, err)
        resp := response_base.NewErrorResponse(errApi)
        return controller.WriteBadRequest(c, resp)
    }

    return controller.WriteSuccessEmptyContent(c)
}

func (controller *BotController) UploadFile(c echo.Context) error {

    formValue, err := c.MultipartForm()
    if err != nil {
        errApi := error_base.New(error_base.ErrorValidDataCode, fmt.Errorf(err.Error()))
        resp := response_base.NewErrorResponse(errApi)
        return controller.WriteBadRequest(c, resp)
    }

    ctx := c.Request().Context()
    if ctx == nil {
        ctx = context.Background()
    }

    rs, err := controller.Service.UploadFiles(ctx, formValue.File["file"])
    if err != nil {
        errApi := error_base.New(error_base.ErrorValidDataCode, fmt.Errorf(err.Error()))
        resp := response_base.NewErrorResponse(errApi)
        return controller.WriteBadRequest(c, resp)
    }

    return controller.WriteSuccess(c, rs)
}
