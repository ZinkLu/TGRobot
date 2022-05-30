package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type Config struct {
	ApiToken string `json:"apiToken"`
	Debug    bool   `json:"debug"`
}

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

func LoadConfig() (*Config, string) {
	dir, _ := os.Getwd()
	config_path := path.Join(dir, "config", "config_secret.json")
	content, err := ioutil.ReadFile(config_path)
	if err != nil {
		panic(err.Error())
	}
	config := &Config{}
	json.Unmarshal(content, config)
	return config, string(content)
}

func LoadDebugConfig() (*Config, string) {
	dir, _ := Dirname()
	config_path := path.Join(dir, "config_secret.json")
	content, err := ioutil.ReadFile(config_path)
	if err != nil {
		panic(err.Error())
	}
	config := &Config{}
	json.Unmarshal(content, config)
	return config, string(content)
}
