package message_handler

import (
	vmshell "github.com/ZinkLu/TGRobot/handlers/message_handler/vmshell"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageHandlerInterface interface {
	Handle(*tgbotapi.Message, *tgbotapi.BotAPI)
}

type MessageHandler struct {
	subHandlers []MessageHandlerInterface
}

func (h MessageHandler) Handle(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message == nil {
		return
	}
	for _, handler := range h.subHandlers {
		go handler.Handle(update.Message, bot)
	}
}

// add a handler if you like
func NewMessageHandler(configString string) *MessageHandler {
	res := make([]MessageHandlerInterface, 0, 0)
	h1, err := vmshell.New(configString)
	if err == nil {
		res = append(res, h1)
	}
	return &MessageHandler{subHandlers: res}
}
