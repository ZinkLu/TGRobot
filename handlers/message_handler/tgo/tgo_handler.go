// https: //github.com/p4gefau1t/trojan-go/blob/master/api/service/client.go

package tgo

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/ZinkLu/TGRobot/config"
	"github.com/ZinkLu/TGRobot/handlers/message_handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const HELP = "查询\"我的流量\"(私信我)"
const COMMAND = "我的流量"
const PASS_QUERY = "输入秘钥（请选择本条消息后点击reply再发送）"
const CONVERT = float32(1024 * 1024)

var CACHE = make(map[int64]string)

const FORMATTER = `USAGE               ⏳ %.2fMB(%.2fGb).
SPEED                ⬆️%.2fkb/s, ⬇️%.2fkb/s
SPEED LIMIT      ⬆️%.2fkb/s, ⬇️%.2fkb/s
ONLINE DEVICE 💻 %d(current) %d(total)

如果想获取剩余流量，可以输入 "服务器流量" 查询
`

type TGoHandler struct {
	Client TrojanServerServiceClient
}

func formatTraffic(response *GetUsersResponse) string {
	status := response.Status
	speed := status.GetSpeedCurrent()
	traffic := status.GetTrafficTotal()
	speedLimit := status.GetSpeedLimit()
	currentIp := status.GetIpCurrent()
	limitIp := status.GetIpLimit()
	if speed == nil || traffic == nil {
		return "发生了一些错误"
	}
	// speed bytes
	dspeed := float32(speed.GetDownloadSpeed()) / 1024.0 // kb/s
	uspeed := float32(speed.GetUploadSpeed()) / 1024.0   // kb/s
	totalTraffic := float32(traffic.GetDownloadTraffic()+traffic.GetUploadTraffic()) * 2.0 / CONVERT
	return fmt.Sprintf(FORMATTER, totalTraffic, totalTraffic/1024.0, uspeed, dspeed, float32(speedLimit.GetDownloadSpeed())/1024.0, float32(speedLimit.GetUploadSpeed())/1024.0, currentIp, limitIp)
}

func (tgo *TGoHandler) GetUserStatus(password string) (*GetUsersResponse, error) {
	if tgo.Client == nil {
		log.Println("tgo.Client has not be init, Call tgo.Init first")
		return nil, fmt.Errorf("no client available")
	}

	user := &User{Password: password}

	log.Println("Call rpc from server")
	// response, err := tgo.Client.GetTraffic(context.TODO(), in)
	client, err := tgo.Client.GetUsers(context.TODO())
	if err != nil {
		log.Printf("call rpc server failed")
		log.Print(err.Error())
		return nil, err
	}
	defer client.CloseSend()

	client.Send(&GetUsersRequest{User: user})
	response, err := client.Recv()

	if err != nil {
		log.Printf("get user failed, %s\n", err)
		return nil, err
	}

	if !response.GetSuccess() {
		log.Printf("get user failed, %s\n", response.GetInfo())
		return response, fmt.Errorf("invalid user")
	}
	log.Printf("Call rpc ends, %s\n", response.GetInfo())

	return response, nil
}

// cache use password
// TODO: add admin password for more actions
func (tgo *TGoHandler) Handle(u *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	u_id := u.Message.From.ID
	msg := u.Message
	var password string
	var ok bool
	// always handler user's reply
	if msg.ReplyToMessage != nil && msg.ReplyToMessage.Text == PASS_QUERY {
		password = msg.Text
		// cache if success;
		status, err := tgo.GetUserStatus(password)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "密码不正确，请重新再试"))
			return
		}
		CACHE[u_id] = password
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, formatTraffic(status)))
		return
	}
	// read from cache
	if password, ok = CACHE[u_id]; ok {
		status, err := tgo.GetUserStatus(password)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "密码不正确，请重新再试"))
		}
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, formatTraffic(status)))
		return
	}

	// or send a password query
	bot.Send(tgbotapi.NewMessage(msg.Chat.ID, PASS_QUERY))
}

// only reply private message
func (tgo *TGoHandler) When(u *tgbotapi.Update) bool {
	if u.Message.Chat.Type != "private" {
		return false
	}
	msg := u.Message

	if strings.Contains(msg.Text, COMMAND) || (msg.ReplyToMessage != nil && msg.ReplyToMessage.Text == PASS_QUERY) {
		return true
	}
	return false
}

func (tgo *TGoHandler) Init(config *config.ConfigUnmarshaler) {
	var c = &Config{}
	var err = config.UnmarshalConfig(c, tgo.Name())
	var conn *grpc.ClientConn

	if err != nil {
		log.Printf("%s handler init failed %s", tgo.Name(), err)
		panic(err)
	}

	if c.UserCert && c.CertPath != "" && c.CertKey != "" {
		clientCert, certErr := tls.LoadX509KeyPair(c.CertPath, c.CertKey)
		if certErr != nil {
			log.Fatalf(certErr.Error())
		}

		var tlsConfig = &tls.Config{
			InsecureSkipVerify: !c.Verify,
			Certificates:       []tls.Certificate{clientCert},
			ServerName:         c.SNI,
		}
		if len(c.CaPaths) > 0 {
			pool := tlsConfig.RootCAs
			for _, path := range c.CaPaths {
				bytes, ioErr := ioutil.ReadFile(path)
				if ioErr != nil {
					continue
				}
				pool.AppendCertsFromPEM(bytes)
			}
		}

		var tlsCredential = credentials.NewTLS(tlsConfig)
		conn, err = grpc.Dial(fmt.Sprintf("%s:%d", c.Addr, c.Port), grpc.WithTransportCredentials(tlsCredential))
	} else {
		conn, err = grpc.Dial(fmt.Sprintf("%s:%d", c.Addr, c.Port), grpc.WithInsecure())
	}
	if err != nil {
		log.Printf("%s handler init failed while connect to gRPC server %s", tgo.Name(), err)
		panic(err)
	}
	tgo.Client = NewTrojanServerServiceClient(conn)

}

func (tgo *TGoHandler) Order() int {
	return 3
}

func (tgo *TGoHandler) Help() string {
	return HELP
}

func (tgo *TGoHandler) Name() string {
	return "Tgo"
}
func init() {
	message_handler.Register(&TGoHandler{})
}
