package send_keyboard_test_handler

import (
	"strings"

	"github.com/ZinkLu/TGRobot/config"
	"github.com/ZinkLu/TGRobot/handlers/message_handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var th = &TestHandler{}

type TestHandler struct{}

// th *TestHandler MessageHandlerInterface
func (th *TestHandler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := update.Message
	b1 := tgbotapi.NewInlineKeyboardButtonData("🔵💊", "blue_pill")
	b2 := tgbotapi.NewInlineKeyboardButtonData("🔴💊", "red_pill")
	markup := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(b1, b2))
	reply := tgbotapi.NewMessage(msg.Chat.ID, "Which you choose")
	reply.ReplyMarkup = markup
	bot.Send(reply)
}

func (th *TestHandler) Init(_ *config.ConfigUnmarshaler) {
}

func (th *TestHandler) When(update *tgbotapi.Update) bool {
	return strings.Contains(update.Message.Text, "choice")
}

func (th *TestHandler) Order() int {
	return 999
}

func (th *TestHandler) Help() string {
	return ""
}

func (th *TestHandler) Name() string {
	return "Test"
}

func init() {
	message_handler.Register(th)
}
