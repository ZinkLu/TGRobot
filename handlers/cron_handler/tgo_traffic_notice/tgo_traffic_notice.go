package tgo_traffic_notice

import (
	"context"
	"fmt"
	"time"

	"github.com/ZinkLu/TGRobot/config"
	"github.com/ZinkLu/TGRobot/handlers/cron_handler"
	"github.com/ZinkLu/TGRobot/handlers/message_handler/tgo"
	"github.com/ZinkLu/TGRobot/pool"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const CONVERT = float32(1024 * 1024)

var th = &tGoTrafficNoticeCronHandler{}

type tGoTrafficNoticeCronHandler struct {
	config       *Config
	pr           *PercentReport
	currentMonth time.Month
}

type Config struct {
	Enabled      bool    `configKey:"enabled"`
	Interval     int     `configKey:"interval"`      // second
	TotalTraffic float32 `configKey:"total_traffic"` // mb
}

func (h *tGoTrafficNoticeCronHandler) Handle(_ *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if h.currentMonth != time.Now().Month() {
		h.currentMonth = time.Now().Month()
		h.pr.Clear()
	}
	tgoH := pool.GetAppHandlerByName[*tgo.TGoHandler]("Tgo")
	if client, err := tgoH.Client.ListUsers(context.TODO(), &tgo.ListUsersRequest{}); err == nil {
		var currentTotal float32 = 0.0
		var maxUsage float32 = 0.0
		for response, err := client.Recv(); err == nil; response, err = client.Recv() {
			traffic := response.GetStatus().TrafficTotal
			userTraffic := float32(traffic.GetDownloadTraffic()+traffic.GetUploadTraffic()) * 2.0 / CONVERT
			if userTraffic > maxUsage {
				maxUsage = userTraffic
			}
			currentTotal = currentTotal + userTraffic
		}
		percent := int(currentTotal / h.config.TotalTraffic * 10)

		// notice every 10 percent
		if percent > 10 {
			return
		}
		if h.pr.NeedReport(percent) {
			for k := range cron_handler.GetSubscribe() {
				bot.Send(
					tgbotapi.NewMessage(k,
						fmt.Sprintf("We have reached %d %% of total traffic( %.02f mb )!;\nmax user's usage is %.02f MB",
							percent*10,
							h.config.TotalTraffic,
							maxUsage,
						)))
			}
		}
	}
}

func (h *tGoTrafficNoticeCronHandler) When(_ *tgbotapi.Update) bool {
	return true
}

func (h *tGoTrafficNoticeCronHandler) Init(config *config.ConfigUnmarshaler) {
	c := &Config{}
	config.UnmarshalConfig(c, h.Name())
	h.config = c
	h.pr = &PercentReport{}
}

func (h *tGoTrafficNoticeCronHandler) Order() int {
	return 0
}

func (h *tGoTrafficNoticeCronHandler) Help() string {
	return ""
}

func (h *tGoTrafficNoticeCronHandler) Name() string {
	return "TgoTrafficNotice"
}

func (h *tGoTrafficNoticeCronHandler) Every() int {
	if !h.config.Enabled || h.config.Interval < 0 {
		return -1
	}
	return h.config.Interval
}
func init() {
	cron_handler.Register(th)
}
