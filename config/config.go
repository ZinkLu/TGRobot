package config

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type HandlerConfig struct {
	MessageHandler, CommandHandler *ConfigUnmarshaler
}

type GlobalConfig struct {
	ApiToken       string `configKey:"apiToken"`
	Debug          bool   `configKey:"debug"`
	Handlers       Config `configKey:"handlers"`
	HandlersConfig *HandlerConfig
}

func LoadTgBotConfig(configPath string) *GlobalConfig {
	config := &GlobalConfig{}
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err.Error())
	}
	var filesInfo = strings.Split(configPath, ".")
	var filesExtension = filesInfo[len(filesInfo)-1]
	configUnmarshal := NewConfigUnmarshaler(content, filesExtension)

	configUnmarshal.UnmarshalConfig(config, "")
	addHandlerConfig(config)
	return config
}

// unmarshal handler config for GlobalConfig object
func addHandlerConfig(config *GlobalConfig) {
	var hc = &HandlerConfig{}
	hs := config.Handlers
	if hs == nil {
		panic("config should have 'handlers' property")
	}

	messageHandler, ok := hs["message_handler"]
	if ok {
		mh, ok := messageHandler.(Config)
		if ok {
			hc.MessageHandler = &ConfigUnmarshaler{mh}
		}
	}
	commandHandler, ok := hs["command_handler"]
	if ok {
		ch, ok := commandHandler.(Config)
		if ok {
			hc.CommandHandler = &ConfigUnmarshaler{ch}
		}
	}
	config.HandlersConfig = hc
}

/*
	for debug usage







*/
func Filename() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filename, nil
}

// Dirname is the __dirname equivalent
func Dirname() (string, error) {
	filename, err := Filename()
	if err != nil {
		return "", err
	}
	return filepath.Dir(filename), nil
}

func LoadDebugConfig() *GlobalConfig {
	dir, _ := Dirname()
	config_path := path.Join(dir, "config_secret.json")
	return LoadTgBotConfig(config_path)
}
