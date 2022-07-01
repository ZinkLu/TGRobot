package handlers

import (
	"github.com/ZinkLu/TGRobot/config"
	"github.com/ZinkLu/TGRobot/pool"

	"github.com/ZinkLu/TGRobot/handlers/inline_keyboard_handler"

	"github.com/ZinkLu/TGRobot/handlers/message_handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	/*
	 for register app handler
	*/
	_ "github.com/ZinkLu/TGRobot/handlers/inline_keyboard_handler/inline_keyboard_test_handler"
	_ "github.com/ZinkLu/TGRobot/handlers/inline_keyboard_handler/vmshell"
	_ "github.com/ZinkLu/TGRobot/handlers/message_handler/send_keyboard_test_handler"
	_ "github.com/ZinkLu/TGRobot/handlers/message_handler/tgo"
	_ "github.com/ZinkLu/TGRobot/handlers/message_handler/vmshell"
	_ "github.com/ZinkLu/TGRobot/handlers/message_handler/yiyan_handler"
)

type TelegramHandlerInterface interface {
	Handle(tgbotapi.Update, *tgbotapi.BotAPI)
}

func GetHandlers(config *config.GlobalConfig) []TelegramHandlerInterface {
	// messageHandler
	message_handler.MESSAGE_HANDLER.Init(config.HandlersConfig.MessageHandler)
	pool.AddAppHandler(message_handler.MESSAGE_HANDLER.AppHandlers...)

	inline_keyboard_handler.INLINE_KEYBOARD_HANDLER.Init(config.HandlersConfig.InlineKeyBoardHandler)
	pool.AddAppHandler(inline_keyboard_handler.INLINE_KEYBOARD_HANDLER.AppHandlers...)

	return []TelegramHandlerInterface{
		message_handler.MESSAGE_HANDLER,
		inline_keyboard_handler.INLINE_KEYBOARD_HANDLER,
	}
}
