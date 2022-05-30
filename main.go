package main

import (
	bot "github.com/ZinkLu/TGRobot/bots"
	config "github.com/ZinkLu/TGRobot/config"
	handlers "github.com/ZinkLu/TGRobot/handlers"
	"github.com/ZinkLu/TGRobot/handlers/message_handler"
	// message_handler
)

func main() {
	config, fullConfig := config.LoadConfig()
	message_handler := message_handler.NewMessageHandler(fullConfig)
	hs := []handlers.TgHandler{message_handler}
	bot.GetBotAndListen(config.ApiToken, config.Debug, hs)
}
