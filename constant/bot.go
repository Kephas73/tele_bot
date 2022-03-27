package constant

const ValueEmpty int = 0

var ReplyMsg = []string{
    "Xin chào %s, đây là hệ thống tự động!",
    "Xin chào %s",
    "Xin chào %s, tôi có thể giúp gì cho bạn?",
    "Xin chào %s, bạn đang tham gia vào group này!",
}

const (
    InitBot      string = "👉 Bot: %s đã bật 😍 😍 😍"
    ReconnectBot string = "👉 Bot: %s kết nối lại thành công 😍 😍 😍"
    ShutdownBot  string = "👉 Bot: %s đã tắt ☹ ☹ ☹"
    IsBot        string = "✅ "
    BotAliveMsg         = "Tôi vẫn còn hoạt động 😜 😜 😜"
    HelpMsg      string = `Bạn có thể sử dụng các lệnh bên dưới:

help - thông tin các lệnh có thể sử dụng
health - kiểm tra bot còn sống hay không
now - ngày giờ hiện tại
name - lấy tên bot
other - tự động trả lời
`
)

const (
    ExpiresThenMinute = 10
)
