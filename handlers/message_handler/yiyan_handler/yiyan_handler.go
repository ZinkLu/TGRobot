package yiyan

import (
	"strings"

	"github.com/ZinkLu/TGRobot/config"
	"github.com/ZinkLu/TGRobot/handlers/message_handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type YiYanHandler struct{}

func (yiyan *YiYanHandler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	m := update.Message
	yy, err := GetYiYan()
	if err == nil {
		bot.Send(tgbotapi.NewMessage(m.Chat.ID, yy.Quote()))
	}
}

func (yiyan *YiYanHandler) Init(_ *config.ConfigUnmarshaler) {
}

func (yiyan *YiYanHandler) When(update *tgbotapi.Update) bool {
	return strings.Contains(update.Message.Text, "一句话")
}

func (yiyan *YiYanHandler) Order() int {
	return 2
}

func (yiyan *YiYanHandler) Help() string {
	return "对我说\"一句话\""
}

func (yiyan *YiYanHandler) Name() string {
	return "一言"
}

func init() {
	message_handler.Register(&YiYanHandler{})
}
