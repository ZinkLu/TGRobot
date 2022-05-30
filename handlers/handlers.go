package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type TgHandler interface {
	Handle(tgbotapi.Update, *tgbotapi.BotAPI)
}
