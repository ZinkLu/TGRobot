package vmshell

import (
	"fmt"
	"testing"

	"github.com/ZinkLu/TGRobot/config"
)

func TestVmShellClient_GetServerInfo(t *testing.T) {
	_, fullConfig := config.LoadDebugConfig()
	// fmt.Printf("%s %s %s\n", config.Username, config.Password, config.Serverid)
	handler, err := New(fullConfig)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	serverInfo, err := handler.client.GetServerInfo(handler.serverId, true)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(serverInfo.GetBandWithStatus())
}
