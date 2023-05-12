package cron_handler

import (
	"github.com/ZinkLu/TGRobot/config"
	"github.com/ZinkLu/TGRobot/handlers/common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jasonlvhit/gocron"
)

var CRON_HANDLER = &CronHandler{started: false, registerChats: make(map[int64]bool), AppHandlers: make([]common.CronHandlerInterface, 0), signal: make(chan *tgbotapi.BotAPI, 1)}
var EMPTY_UPDATE = &tgbotapi.Update{}

type CronHandler struct {
	AppHandlers   []common.CronHandlerInterface
	started       bool
	registerChats map[int64]bool
	switches      chan bool
	signal        chan *tgbotapi.BotAPI
}

func (h *CronHandler) helper() {
	for {
		bot := <-h.signal
		if bot != nil && !h.started {
			h.started = true
			for _, v := range CRON_HANDLER.AppHandlers {
				if v.Every() < 0 {
					continue
				}
				gocron.Every(uint64(v.Every())).Seconds().Do(v.Handle, EMPTY_UPDATE, bot)
			}
			h.switches = gocron.Start()
		}
		if bot == nil && len(h.registerChats) == 0 {
			h.started = false
			gocron.Clear()
			if h.switches != nil {
				h.switches <- false
			}
		}
	}
}

func (h *CronHandler) handleSubscribe(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message == nil {
		return
	}
	chatId := update.Message.Chat.ID
	switch update.Message.Text {
	case "subscribe":
		h.signal <- bot
		_, ok := h.registerChats[chatId]
		if ok {
			bot.Send(tgbotapi.NewMessage(chatId, "you had subscribed"))
		} else {
			h.registerChats[chatId] = true
			bot.Send(tgbotapi.NewMessage(chatId, "subscribed"))
		}
	case "unsubscribe":
		h.signal <- nil
		_, ok := h.registerChats[chatId]
		if ok {
			delete(h.registerChats, chatId)
			bot.Send(tgbotapi.NewMessage(chatId, "unsubscribed"))

		} else {
			bot.Send(tgbotapi.NewMessage(chatId, "you are not subscribed"))
		}
	}
}

func (h *CronHandler) Handle(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	h.handleSubscribe(&update, bot)
}

func (mh *CronHandler) Init(configUnmarshaler *config.ConfigUnmarshaler) {
	for _, mh := range mh.AppHandlers {
		mh.Init(configUnmarshaler)
	}
	go mh.helper()
}

// inline keyboard doesn't have any
func (mh *CronHandler) Help() string {
	panic("no help information needed")
}

// call Register to enable a handler
func Register(h common.CronHandlerInterface) {
	CRON_HANDLER.AppHandlers = append(CRON_HANDLER.AppHandlers, h)
}

func GetSubscribe() map[int64]bool {
	return CRON_HANDLER.registerChats
}
