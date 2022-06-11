package vmshell

import (
	"strings"

	"github.com/ZinkLu/TGRobot/config"
	"github.com/ZinkLu/TGRobot/handlers/message_handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

const errorMessage = "服务器出了点问题☢️"
const helpMessage = "查\"流量\"，查询\"服务器状态\""

const (
	USAGE = "流量"
	INFO  = "服务器"
)

var replyMessage = [2]string{USAGE, INFO}

type VmShellHandler struct {
	client   *vmShellClient
	serverId string
}

func (v *VmShellHandler) getCorrectMessage(msg *tgbotapi.Message) string {
	msg_string := msg.Text

	msgType := getMessageType(msg_string)

	switch msgType {
	case "":
		return helpMessage
	case USAGE:
		si, err := v.client.GetServerInfo(v.serverId, true)
		if err != nil {
			log.Info(err)
			return errorMessage
		}
		return si.GetBandWithStatus()
	case INFO:
		si, err := v.client.GetServerInfo(v.serverId, true)
		if err != nil {
			log.Info(err)
			return errorMessage
		}
		return si.GetServerStatus()
	default:
		return "你已经进入了异次元，你是怎么进来的？"
	}

}

func (vh *VmShellHandler) Handle(msg *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	var message = ""
	message = vh.getCorrectMessage(msg)
	sendMessage := tgbotapi.NewMessage(msg.Chat.ID, message)
	bot.Send(sendMessage)
}

// init Client
func (vh *VmShellHandler) Init(config *config.ConfigUnmarshaler) {
	c := &Config{}
	config.UnmarshalConfig(c, vh.Name())
	vh.client = newClient(c.Username, c.Password)
	vh.serverId = c.ServerId
}

func (vh *VmShellHandler) When(msg *tgbotapi.Message) bool {
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
