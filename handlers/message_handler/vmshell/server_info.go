package vmshell

import (
	"fmt"
	"strconv"
	"strings"
)

type ServerInfo struct {
	Bandwidthcolor       string   `json:"bandwidthcolor"`
	Bandwidthfree        string   `json:"bandwidthfree"`
	Bandwidthpercent     string   `json:"bandwidthpercent"`
	Bandwidthtotal       string   `json:"bandwidthtotal"`
	Bandwidthused        string   `json:"bandwidthused"`
	Clientkeyautherror   int      `json:"clientkeyautherror"`
	Connaddr             string   `json:"connaddr"`
	Cpus                 string   `json:"cpus"`
	Disk                 string   `json:"disk"`
	Displaybandwidthbar  int      `json:"displaybandwidthbar"`
	Displayboot          int      `json:"displayboot"`
	Displayclientkeyauth int      `json:"displayclientkeyauth"`
	Displayconsole       int      `json:"displayconsole"`
	Displaygraphs        int      `json:"displaygraphs"`
	Displayhddbar        int      `json:"displayhddbar"`
	Displayhostname      int      `json:"displayhostname"`
	Displayhtml5console  int      `json:"displayhtml5console"`
	Displayips           int      `json:"displayips"`
	Displayloadgraph     int      `json:"displayloadgraph"`
	Displaymemorybar     int      `json:"displaymemorybar"`
	Displaymemorygraph   int      `json:"displaymemorygraph"`
	Displaypanelbutton   int      `json:"displaypanelbutton"`
	Displayreboot        int      `json:"displayreboot"`
	Displayrebuild       int      `json:"displayrebuild"`
	Displayrootpassword  int      `json:"displayrootpassword"`
	Displayshutdown      int      `json:"displayshutdown"`
	Displaystatus        string   `json:"displaystatus"`
	Displaytrafficgraph  int      `json:"displaytrafficgraph"`
	Displayvnc           int      `json:"displayvnc"`
	Displayvncpassword   int      `json:"displayvncpassword"`
	Firstport            string   `json:"firstport"`
	Hostname             string   `json:"hostname"`
	Ipcsv                []string `json:"ipcsv"`
	Lastport             string   `json:"lastport"`
	Mac                  string   `json:"mac"`
	Mainip               string   `json:"mainip"`
	Memory               string   `json:"memory"`
	Mode                 string   `json:"node"`
	Sshport              string   `json:"sshport"`
	State                string   `json:"state"`
	Swap                 string   `json:"swap"`
	Template             string   `json:"template"`
	Trafficgraphurl      string   `json:"trafficgraphurl"`
	Type                 string   `json:"type"`
}

func getPercentBar(percent int) string {
	percent = percent / 2
	sb := strings.Builder{}
	for i := 0; i <= 50; i++ {
		if i <= percent {
			sb.WriteString(">")
		} else {
			sb.WriteString("=")
		}
	}
	return sb.String()
}

func (si *ServerInfo) GetBandWithStatus() string {
	var bar string = ""
	i, err := strconv.ParseInt(si.Bandwidthpercent, 10, 64)
	if err != nil {
		bar = ""
	}
	bar = getPercentBar(int(i))
	return fmt.Sprintf("total: %s, used: %s, free: %s \nuse %s percent : [%s]", si.Bandwidthtotal, si.Bandwidthused, si.Bandwidthfree, si.Bandwidthpercent, bar)
}

func (si *ServerInfo) GetServerStatus() string {
	formatter := "ServerIP: %s\nDisk: %s\nMemory: %s\nCpus: %s\n"
	return fmt.Sprintf(formatter, si.Connaddr, si.Disk, si.Memory, si.Cpus)
}
