package service

import (
    "botTele/constant"
    "botTele/infrastructure/bot"
    "botTele/infrastructure/logger"
    "botTele/model"
    "bytes"
    "encoding/json"
    "fmt"
    "time"
)

type IBotService interface {
    SendChat(data model.RawData) error
    SendChatShutdown() error
    AutoReply()
}

type BotService struct {
    Bot     *bot.TelegramBot
    Timeout time.Duration
}

func NewBotService(timeout time.Duration) IBotService {
    return &BotService{
        Bot:     bot.NewBotTele(),
        Timeout: timeout,
    }
}

func (bot *BotService) SendChat(data model.RawData) error {

    go func() {
        if len(data.Text) != constant.ValueEmpty {
            err := bot.Bot.SendMessageForApiTele(data.Text)
            if err != nil {
                logger.Error("BotService::SendChat: - Send chat error: %v", err)
            }
        }
    }()

    go func() {
        if data.Object != nil {
            b, _ := json.Marshal(data.Object)
            err := bot.Bot.SendMessageForApiTele(bytes.NewBuffer(b))
            if err != nil {
                logger.Error("BotService::SendChat: - Send chat error: %v", err)
            }
        }
    }()

    return nil
}

func (bot *BotService) AutoReply() {
    bot.Bot.AutoReply()
}

func (bot *BotService) SendChatShutdown() error {
    err := bot.Bot.SendChat(fmt.Sprintf(constant.ShutdownBot, bot.Bot.Bot.Self.FirstName+" "+bot.Bot.Bot.Self.LastName))
    if err != nil {
        logger.Error("BotService::SendChat: - Send chat error: %v", err)
        return err
    }
    return nil
}
