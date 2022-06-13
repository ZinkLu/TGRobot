package vmshell

import (
	"strings"

	"github.com/ZinkLu/TGRobot/config"
	"github.com/ZinkLu/TGRobot/handlers/inline_keyboard_handler"
	vm_message "github.com/ZinkLu/TGRobot/handlers/message_handler/vmshell"
	"github.com/ZinkLu/TGRobot/pool"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var th = &vmShellInlineKeyBoardHandler{}

type vmShellInlineKeyBoardHandler struct{}

func (th *vmShellInlineKeyBoardHandler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	callback := update.CallbackQuery
	reply := tgbotapi.NewMessage(callback.Message.Chat.ID, "")
	var err error

	commands := strings.Split(callback.Data, "_")
	command, s_id := commands[1], commands[2]
	switch command {
	case vm_message.USAGE:
		message_handler := pool.GetAppName[*vm_message.VmShellHandler]("vmShell")
		if si, err := message_handler.Client.GetServerInfo(s_id, true); err == nil {
			reply.Text = si.GetBandWithStatus()
		}
	case vm_message.INFO:
		message_handler := pool.GetAppName[*vm_message.VmShellHandler]("vmShell")
		if si, err := message_handler.Client.GetServerInfo(s_id, true); err == nil {
			reply.Text = si.GetServerStatus()
		}
	default:
		reply.Text = "你进入了异次元"
	}

	if err != nil {
		reply.Text = vm_message.ErrorMessage
	}

	bot.Send(reply)
}

func (th *vmShellInlineKeyBoardHandler) Init(_ *config.ConfigUnmarshaler) {
}

func (th *vmShellInlineKeyBoardHandler) When(update *tgbotapi.Update) bool {
	callback := update.CallbackQuery
	return strings.Contains(callback.Data, vm_message.PREFIX)
}

func (th *vmShellInlineKeyBoardHandler) Order() int {
	return 0
}

func (th *vmShellInlineKeyBoardHandler) Help() string {
	panic("no help information needed")
}

func (th *vmShellInlineKeyBoardHandler) Name() string {
	return "vmShellInlineKeyBoardHandler"
}

func init() {
	inline_keyboard_handler.Register(th)
}