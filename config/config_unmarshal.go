package config

/*
	provider a standard unmarshaler to
	different format of config files, such
	as yaml, json, etc.
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"
)

type Config = map[string]interface{} // consider config as key: value structure

// for config unmarshal.
type ConfigUnmarshaler struct {
	config Config
}

/*	UnmarshalConfig

	unmarshal file to a config object;

	an optional filed `configKey` tag can be set to marshal

	type config struct {
		Name string `configKey: "name"`
	}

	@param name:
		if you want to unmarshal a sub config object,
		pass the sub config key as name param.
*/
func (ah ConfigUnmarshaler) UnmarshalConfig(config interface{}, name string) error {
	// which means to unmarshal all field
	var appValue Config
	// var ok bool

	if name != "" {
		value, ok := ah.config[name]
		if !ok {
			return errors.New("not config name: " + name)
		}
		appValue, ok = value.(Config)
		if !ok {
			return errors.New("app config should be a map or a object")
		}
	} else {
		appValue = ah.config
	}

	c := reflect.Indirect(reflect.ValueOf(config))

	getType := c.Type()
	// getValue := reflect.ValueOf(c)

	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i) // real key which should be get from value
		key := field.Tag.Get("configKey")
		if key == "" {
			key = field.Name
		} // use filed name instead
		kValue, ok := appValue[key]
		if ok {
			v := reflect.ValueOf(kValue)
			f := c.Field(i)
			if f.CanSet() {
				f.Set(v) // wrong type should panic
			}
		}
	}
	return nil
}

func NewConfigUnmarshaler(content []byte, ext string) *ConfigUnmarshaler {
	var c = Config{}

	err := fileUnmarshal(content, ext, &c)

	if err != nil {
		panic(err.Error())
	}
	return &ConfigUnmarshaler{config: c}
}

func fileUnmarshal(content []byte, ext string, obj interface{}) error {
	switch ext {
	case "json":
		return json.Unmarshal(content, obj)
	case "yaml", "yml":
		return yaml.Unmarshal(content, obj)
	default:
		return fmt.Errorf("no Unmarshal found")
	}
}
