// app handler is a true handler which handler different message

package common

import (
	"github.com/ZinkLu/TGRobot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AppHandlerInterface interface {
	Handle(*tgbotapi.Update, *tgbotapi.BotAPI)
	When(*tgbotapi.Update) bool // for Chain of Responsibility
	Init(*config.ConfigUnmarshaler)
	Order() int // for Chain of Responsibility, less is higher
	Help() string
	Name() string
}

type CronHandlerInterface interface {
	AppHandlerInterface
	Every() int // disabled if return value <0
}
