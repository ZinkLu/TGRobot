package main

import (
	"flag"

	// message_handler
	bot "github.com/ZinkLu/TGRobot/bots"
	"github.com/ZinkLu/TGRobot/config"
	"github.com/ZinkLu/TGRobot/handlers"
	log "github.com/sirupsen/logrus"
)

func main() {
	var configPath = flag.String("c", "", "config file path")
	flag.Parse()

	if *configPath == "" {
		log.Error("config file path should be set")
		return
	}

	config := config.LoadTgBotConfig(*configPath)

	hs := handlers.GetHandlers(config)
	bot.GetBotAndListen(config.ApiToken, config.Debug, hs)
}
