package bot

import (
    "botTele/constant"
    "botTele/infrastructure/logger"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "net/http"
    "time"
)

var tr = &http.Transport{
    MaxIdleConns:          50,
    MaxIdleConnsPerHost:   50,
    ResponseHeaderTimeout: 30 * time.Second,
}

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

func (botTele *TelegramBot) SendMessageForApiTele(sendMsg interface{}) error {

    sendMsgUri := fmt.Sprintf(botTele.SendMessageUri, botTele.Token, botTele.ChatId[0], IsBotMsg(fmt.Sprint(sendMsg)))
    req, err := http.NewRequest("POST", sendMsgUri, nil)
    if err != nil {
        logger.Error("TelegramBot::SendMessageForApiTele - Can not create new request error: %v", err)
        return err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{
        Timeout:   time.Second * 10,
        Transport: tr,
    }

    response, err := client.Do(req)
    if err != nil {
        logger.Error("TelegramBot::SendMessageForApiTele - Can not exec request error: v", err)
        return err
    }

    if response.StatusCode != 200 {
        logger.Error("TelegramBot::SendMessageForApiTele - Send msg for api error code: %d", response.StatusCode)
        return fmt.Errorf("send msg for api error code: %d", response.StatusCode)
    }

    return nil
}