<p align="center">
 <img width="100px" src="https://res.cloudinary.com/anuraghazra/image/upload/v1594908242/logo_ccswme.svg" align="center" alt="GitHub Readme Stats" />
 <h2 align="center">Bot telegram</h2>
 <p align="center">Get dynamically generated GitHub stats on your readmes!</p>
</p>
  <p align="center">
    <a href="https://a.paddle.com/v2/click/16413/119403?link=1227">
      <img src="https://img.shields.io/badge/Supported%20by-VSCode%20Power%20User%20%E2%86%92-gray.svg?colorA=655BE1&colorB=4F44D6&style=for-the-badge"/>
    </a>
  </p>
</p>


# Features
- [Create bot and get chat id](#create-bot-and-get-chat-id)
- [Send api bot](#send-api-bot)
- [Auto reply](#auto-reply)

### Create bot and get chat id

Link create bot: https://core.telegram.org/bots#creating-a-new-bot

Link get chat id: https://www.alphr.com/find-chat-id-telegram/

### Send api bot

POST: /bot/send-chat
```json
{
    "text":"",
    "object":{}
}
```

### Auto reply
Auto reply messages

### Send message for api telegram
```
Set config: 
    "TelegramApi": {
        "SendMessageUri": "https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s",
        "Token": "xxxxxxxxxxx",
        "ChatId": [
            -xxxxxxx
        ],
        "TimeDelay": 2,
        "ExpirationTime": 1
    }
    
Note: 
    Api support: https://api.telegram.org/bot<token>/sendMessage?chat_id=<chat_id>&text=<msg>    
```
Demo: func SendMessageForApiTele()

---

Contributions are welcome! <3