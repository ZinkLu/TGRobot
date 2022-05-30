package vmshell

import (
	"encoding/json"

	gjson "github.com/tidwall/gjson"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ServerId string `json:"serverid"`
}

func fromJsonToConfig(jsonString string) *Config {
	config := &Config{}
	value := gjson.Get(jsonString, "handlers.message_handler.vmshell")
	// fmt.Println(value)
	json.Unmarshal([]byte(value.String()), config)
	return config
}
