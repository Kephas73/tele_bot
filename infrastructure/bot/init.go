package bot

import (
    "botTele/constant"
    "botTele/infrastructure/logger"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/spf13/viper"
    "math/rand"
    "strings"
    "time"
)

type Config struct {
    Token          string  `json:"token,omitempty"`
    ExpirationTime int     `json:"expiration_time,omitempty"`
    ChatId         []int64 `json:"chat_id,omitempty"`
    TimeDelay      int     `json:"time_delay,omitempty"`
}

var (
    config *Config
)

type TelegramBot struct {
    Bot               *tgbotapi.BotAPI
    BotExpirationTime int64
    ExpirationTime    int
    Token             string
    ChatId            []int64
    TimeDelay         int
    ReplyMsg          []string
}

var teleBot *TelegramBot

func NewBotTele() *TelegramBot {
    if teleBot == nil {
        configKey := "TelegramApi"
        config = &Config{}

        if err := viper.UnmarshalKey(configKey, config); err != nil {
            err = fmt.Errorf("not found config name with env %q for amqp with error: %+v", configKey, err)
            panic(err)
        }

        if config.ExpirationTime <= 0 {
            config.ExpirationTime = constant.ExpiresThenMinute
        }

        bot, err := tgbotapi.NewBotAPI(config.Token)
        if err != nil {
            logger.Info("Failed bot error: %v", err)
            panic(err)
        }
        bot.Debug = true
        logger.Info("Success bot: Authorized on account %s", bot.Self.UserName)

        teleBot = &TelegramBot{
            Bot:               bot,
            BotExpirationTime: time.Now().Add(time.Minute * time.Duration(config.ExpirationTime)).Unix(),
            Token:             config.Token,
            ChatId:            config.ChatId,
            TimeDelay:         config.TimeDelay,
        }

        teleBot.SendChat(fmt.Sprintf(constant.InitBot, teleBot.Bot.Self.FirstName+" "+teleBot.Bot.Self.LastName))
    }

    return teleBot
}

func (botTele *TelegramBot) AutoReply() {
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    updates := botTele.Bot.GetUpdatesChan(u)

    // Loop through each update.
    for update := range updates {
        // Check if we've gotten a message update.
        if update.Message != nil && update.Message.From.IsBot == false {
            // Construct a new message from the given chat ID and containing
            // the text that we received.
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
            // If the message was open, add a copy of our numeric keyboard.
            switch strings.ToLower(update.Message.Text) {
            case "health":
                msg.Text = IsBotMsg(constant.BotAliveMsg)
            case "now":
                msg.Text = IsBotMsg(fmt.Sprintf("Now: %s", time.Now().Format("2006-01-02 15-04-05")))
            case "name":
                msg.Text = IsBotMsg(fmt.Sprintf("Tên tôi là %s", botTele.Bot.Self.FirstName+" "+botTele.Bot.Self.LastName))
            case "help":
                msg.Text = IsBotMsg(constant.HelpMsg)
            default:
                msg.ReplyToMessageID = update.Message.MessageID
                rand.Seed(time.Now().UnixNano())
                msg.Text = IsBotMsg(fmt.Sprintf(constant.ReplyMsg[rand.Intn(len(constant.ReplyMsg))],
                    update.Message.From.FirstName+" "+update.Message.From.LastName))
            }

            // Send the message.
            if _, err := botTele.Bot.Send(msg); err != nil {
                logger.Error("TelegramBot::SendChat: Bot send chat error: %v", err)
            }
        }
    }
}

func IsBotMsg(msg string) string {
    return fmt.Sprintf("%s %s", constant.IsBot, msg)
}

func (botTele *TelegramBot) ReconnectBotTele() {
    if botTele.BotExpirationTime < time.Now().Unix() {
        bot, err := tgbotapi.NewBotAPI(botTele.Token)
        if err != nil {
            logger.Info("Failed bot error: %v", err)
            panic(err)
        }
        bot.Debug = true
        logger.Info("Success reconnect bot: Authorized on account %s", bot.Self.UserName)

        botTele.Bot = bot
        botTele.BotExpirationTime = time.Now().Add(time.Minute * time.Duration(botTele.ExpirationTime)).Unix()

        //botTele.SendChat(fmt.Sprintf(constant.ReconnectBot, teleBot.Bot.Self.FirstName+" "+teleBot.Bot.Self.LastName))
    }
}
