package inline_keyboard_test_handler

import (
	"strings"

	"github.com/ZinkLu/TGRobot/config"
	"github.com/ZinkLu/TGRobot/handlers/inline_keyboard_handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var th = &TestHandler{}

type TestHandler struct{}

func (th *TestHandler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	callback := update.CallbackQuery
	reply := tgbotapi.NewMessage(callback.Message.Chat.ID, "")
	switch callback.Data {
	case "red_pill":
		reply.Text = "You are in matrix, wake up"
	default:
		reply.Text = "nothing to complain,I love my live"
	}

	bot.Send(reply)
}

func (th *TestHandler) Init(_ *config.ConfigUnmarshaler) {
}

func (th *TestHandler) When(update *tgbotapi.Update) bool {
	callback := update.CallbackQuery
	return strings.Contains(callback.Data, "_pill")
}

func (th *TestHandler) Order() int {
	return 999
}

func (th *TestHandler) Help() string {
	panic("no help information needed")
}

func (th *TestHandler) Name() string {
	return "TestHandler"
}

func init() {
	inline_keyboard_handler.Register(th)
}
