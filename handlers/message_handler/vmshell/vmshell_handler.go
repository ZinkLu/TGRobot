package vmshell

import (
	"fmt"
	"strings"

	"github.com/ZinkLu/TGRobot/config"
	"github.com/ZinkLu/TGRobot/handlers/message_handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type serverInfoList = [][2]string

const ErrorMessage = "服务器出了点问题☢️"
const PREFIX = "vmServerId_"

const helpMessage = "查\"流量\"，查\"信息\""
const (
	USAGE = "流量"
	INFO  = "信息"
)

var replyMessage = [2]string{USAGE, INFO}

type VmShellHandler struct {
	Client    *vmShellClient
	ServerIds []string
}

func (v *VmShellHandler) getInlineKeyboardsMessage(msg *tgbotapi.Message) (serverInfoList, error) {
	result := make(serverInfoList, 0)
	var err error
	msg_string := msg.Text
	msgType := getMessageType(msg_string)
	if msgType != "" {
		result, err = v.getKeyBoardRowsInfo(msgType)
		if err != nil {
			err = fmt.Errorf(ErrorMessage)
		}
	}

	return result, err
}

/*
	return [][2]string

	[2]string means [keyboard text, keyboard, data]
*/
func (v *VmShellHandler) getKeyBoardRowsInfo(command string) (serverInfoList, error) {
	res := make(serverInfoList, len(v.ServerIds))
	infos, err := v.Client.GetServersInfo(v.ServerIds, true)
	if err != nil {
		return res, err
	}

	for idx, info := range infos {
		t := [2]string{
			info.GetOneLineInfo(),                     // text
			PREFIX + command + "_" + v.ServerIds[idx], // data
		}
		res[idx] = t
	}
	return res, nil
}

func (vh *VmShellHandler) Handle(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := update.Message
	message, err := vh.getInlineKeyboardsMessage(msg)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, err.Error()))
		return
	}

	buttons := make([]tgbotapi.InlineKeyboardButton, len(message))
	for idx, m := range message {
		text, data := m[0], m[1]
		buttons[idx] = tgbotapi.NewInlineKeyboardButtonData(text, data)
	}
	if len(buttons) == 0 {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "未找到服务器信息"))
		return
	}

	row := tgbotapi.NewInlineKeyboardRow(buttons...)
	markup := tgbotapi.NewInlineKeyboardMarkup(row)

	tgMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "请问需要查询那台服务器：")

	tgMessage.ReplyMarkup = markup
	bot.Send(tgMessage)
}

// init Client
func (vh *VmShellHandler) Init(config *config.ConfigUnmarshaler) {
	c := &Config{}
	config.UnmarshalConfig(c, vh.Name())
	vh.Client = newClient(c.Username, c.Password)
	vh.ServerIds = c.ServerIds
}

func (vh *VmShellHandler) When(update *tgbotapi.Update) bool {
	msg := update.Message
	for _, v := range replyMessage {
		if strings.Contains(msg.Text, v) {
			return true
		}
	}
	return false
}

func (vh *VmShellHandler) Order() int {
	return 1
}

func (vh *VmShellHandler) Help() string {
	return helpMessage
}

func (vh *VmShellHandler) Name() string {
	return "vmShell"
}

func getMessageType(msg string) string {
	if strings.Contains(msg, USAGE) {
		return USAGE
	} else if strings.Contains(msg, INFO) {
		return INFO
	} else {
		return ""
	}
}

func init() {
	message_handler.Register(&VmShellHandler{})
}
