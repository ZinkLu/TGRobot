package bot

import (
	handler "github.com/ZinkLu/TGRobot/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetBot(apiKey string, debug bool) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		panic(err)
	}

	bot.Debug = debug

	return bot
}

func GetBotAndListen(apiKey string, debug bool, handlers []handler.TelegramHandlerInterface) {
	bot := GetBot(apiKey, debug)

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		for _, handler := range handlers {
			go handler.Handle(update, bot) // Visitor
		}
	}
}
