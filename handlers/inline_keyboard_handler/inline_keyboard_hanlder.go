package inline_keyboard_handler

import (
	"sort"

	"github.com/ZinkLu/TGRobot/config"
	"github.com/ZinkLu/TGRobot/handlers/common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var INLINE_KEYBOARD_HANDLER = &InlineKeyBoardHandler{}

type InlineKeyBoardHandler struct {
	AppHandlers []common.AppHandlerInterface
}

func canReplay(update *tgbotapi.Update) bool {
	return update.CallbackQuery != nil
}

// Handle(tgbotapi.Update, *tgbotapi.BotAPI)
func (h *InlineKeyBoardHandler) Handle(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if !canReplay(&update) {
		return
	}
	for _, handler := range h.AppHandlers {
		if handler.When(&update) {
			handler.Handle(&update, bot)
			return
		}
	}
	// otherwise return help information
	// bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, h.Help()))
}

func (mh *InlineKeyBoardHandler) Init(configUnmarshaler *config.ConfigUnmarshaler) {
	for _, mh := range mh.AppHandlers {
		mh.Init(configUnmarshaler)
	}
}

// inline keyboard doesn't have any
func (mh *InlineKeyBoardHandler) Help() string {
	panic("no help information needed")
}

// call Register to enable a handler
func Register(h common.AppHandlerInterface) {
	INLINE_KEYBOARD_HANDLER.AppHandlers = append(INLINE_KEYBOARD_HANDLER.AppHandlers, h)
	sort.Slice(INLINE_KEYBOARD_HANDLER.AppHandlers, func(i, j int) bool {
		return INLINE_KEYBOARD_HANDLER.AppHandlers[i].Order() < INLINE_KEYBOARD_HANDLER.AppHandlers[j].Order()
	})
}
