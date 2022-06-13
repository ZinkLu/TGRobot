package message_handler

import (
	"fmt"
	"sort"
	"strings"

	config "github.com/ZinkLu/TGRobot/config"
	"github.com/ZinkLu/TGRobot/handlers/common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var MESSAGE_HANDLER = &MessageHandler{}

type MessageHandler struct {
	AppHandlers []common.AppHandlerInterface
}

/*
MessageHandler can handler:
	direct message
	group message with @bot
*/

func CanHandler(msg *tgbotapi.Message, bot *tgbotapi.BotAPI) bool {
	botName := bot.Self.UserName
	canReplay := false
	// 群内的 @ 消息
	if msg.Chat.Type == "group" && strings.Contains(msg.Text, botName) {
		canReplay = true
	}
	// 单聊的消息

	if msg.Chat.Type != "group" && msg.Text != "" {
		canReplay = true
	}
	return canReplay
}

// Handle(tgbotapi.Update, *tgbotapi.BotAPI)
func (h *MessageHandler) Handle(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	// not a message or can not handler
	if update.Message == nil || !CanHandler(update.Message, bot) {
		return
	}

	for _, handler := range h.AppHandlers {
		if handler.When(&update) {
			handler.Handle(&update, bot)
			return
		}
	}

	// otherwise return help information
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, h.Help()))

}

func (mh *MessageHandler) Init(configUnmarshaler *config.ConfigUnmarshaler) {
	for _, mh := range mh.AppHandlers {
		mh.Init(configUnmarshaler)
	}
}

/*
	1. HandlerName:
		HelPInfo
	2. HandlerName:
		HelpInfo
*/
func (mh *MessageHandler) Help() string {
	if len(mh.AppHandlers) == 0 {
		return "我什么也做不了..."
	}
	const HELPER = "我可以：\n"

	sb := strings.Builder{}

	for _, ah := range mh.AppHandlers {
		s := fmt.Sprintf("    %s\n", ah.Help())
		sb.WriteString(s)
	}

	return HELPER + sb.String()
}

// call Register to enable a handler
func Register(h common.AppHandlerInterface) {
	MESSAGE_HANDLER.AppHandlers = append(MESSAGE_HANDLER.AppHandlers, h)
	sort.Slice(MESSAGE_HANDLER.AppHandlers, func(i, j int) bool {
		return MESSAGE_HANDLER.AppHandlers[i].Order() < MESSAGE_HANDLER.AppHandlers[j].Order()
	})
}
