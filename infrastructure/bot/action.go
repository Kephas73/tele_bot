package bot

import (
    "botTele/constant"
    "botTele/infrastructure/logger"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "time"
)

func (botTele *TelegramBot) SendChat(sendMsg interface{}) error {
    botTele.ReconnectBotTele()
    var msg tgbotapi.MessageConfig

    if len(botTele.ChatId) != constant.ValueEmpty {
        for _, v := range botTele.ChatId {
            msg = tgbotapi.NewMessage(v, IsBotMsg(fmt.Sprint(sendMsg)))
            if _, err := botTele.Bot.Send(msg); err != nil {
                logger.Error("TelegramBot::SendChat: Bot send chat error: %v", err)
                return err
            }
            time.Sleep(time.Duration(botTele.TimeDelay) * time.Second)
        }
    }
    return nil
}

func (botTele *TelegramBot) Status() error {
    botTele.ReconnectBotTele()
    var msg tgbotapi.MessageConfig

    if len(botTele.ChatId) != constant.ValueEmpty {
        for _, v := range botTele.ChatId {
            msg = tgbotapi.NewMessage(v, IsBotMsg(constant.BotAliveMsg))
            if _, err := botTele.Bot.Send(msg); err != nil {
                logger.Error("TelegramBot::SendChat: Bot send chat error: %v", err)
                return err
            }
        }
    }
    return nil
}
