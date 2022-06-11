package handlers

import (
	"github.com/ZinkLu/TGRobot/config"
	"github.com/ZinkLu/TGRobot/handlers/message_handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramHandlerInterface interface {
	Handle(tgbotapi.Update, *tgbotapi.BotAPI)
}

func GetHandlers(config *config.GlobalConfig) []TelegramHandlerInterface {
	// messageHandler
	message_handler.MESSAGE_HANDLER.Init(config.HandlersConfig.MessageHandler)

	return []TelegramHandlerInterface{
		message_handler.MESSAGE_HANDLER,
	}
}
