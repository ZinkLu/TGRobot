package vmshell

import (
	"fmt"
	"log"
	"strings"

	yiyan "github.com/ZinkLu/TGRobot/common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const errorMessage = "服务器除了点问题☢️"

type VmShellHandler struct {
	client   *vmShellClient
	serverId string
}

// 只有在群里面 @机器人 或者直接单聊机器人的可以回复
func CanReply(msg *tgbotapi.Message) bool {
	canReplay := false

	// 群内的 @ 消息
	if msg.Chat.Type == "group" && strings.Contains(msg.Text, "@vmshell_network_manager_bot") {
		canReplay = true
	}
	// 单聊的消息

	if msg.Chat.Type != "group" && msg.Text != "" {
		canReplay = true
	}

	return canReplay
}

func (v *VmShellHandler) Handle(msg *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	// 判断是否要发消息
	if !CanReply(msg) {
		return
	}
	var message = ""
	message = v.GetCorrectMessage(msg)
	sendMessage := tgbotapi.NewMessage(msg.Chat.ID, message)
	bot.Send(sendMessage)
}

func (v *VmShellHandler) GetCorrectMessage(msg *tgbotapi.Message) string {
	msg_string := msg.Text

	msgType := GetMessageType(msg_string)

	switch msgType {
	case 0:
		return "你是要查询？流量? 服务器状态? ip地址？或者念一句话？"
	case 1:
		si, err := v.client.GetServerInfo(v.serverId, true)
		if err != nil {
			return errorMessage
		}
		return si.GetBandWithStatus()
	case 2:
		si, err := v.client.GetServerInfo(v.serverId, true)
		if err != nil {
			return errorMessage
		}
		return si.GetServerStatus()
	case 3:
		yiyan, err := yiyan.GetYiYan()
		if err != nil {
			log.Println(err)
			return "em，我似乎也说不出什么话了..."
		}
		return fmt.Sprintf("读读下面这句话吧:\n\n   %s", yiyan.Quote())
	default:
		return "你已经进入了异次元，你是怎么进来的？"
	}

}

func GetMessageType(msg string) int {
	if strings.Contains(msg, "流量") {
		return 1
	} else if strings.Contains(msg, "服务器") {
		return 2
	} else if strings.Contains(msg, "一句话") {
		return 3
	} else {
		return 0
	}
}

func New(jsonString string) (*VmShellHandler, error) {
	config := fromJsonToConfig(jsonString)
	if config.Username == "" {
		return nil, fmt.Errorf("can't construct Config from json %s", jsonString)
	}
	return &VmShellHandler{client: newClient(config.Username, config.Password), serverId: config.ServerId}, nil
}
