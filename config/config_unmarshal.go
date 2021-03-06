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

func processNumber(value reflect.Value, convertTo reflect.Kind) interface{} {
	cType := value.Type().Kind()
	switch cType {
	case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := value.Float()
		switch convertTo {
		case reflect.Float32:
			return float32(v)
		case reflect.Float64:
			return float64(v)
		case reflect.Int:
			return int(v)
		case reflect.Int8:
			return int8(int(v))
		case reflect.Int16:
			return int16(int(v))
		case reflect.Int32:
			return int32(int(v))
		case reflect.Int64:
			return int64(int(v))
		default:
			panic(fmt.Sprintf("can't convert %s type to Number", value.Type().Name()))
		}
	default:
		panic(fmt.Sprintf("can't convert %s type to Number", value.Type().Name()))
	}
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

	confObjValue := reflect.Indirect(reflect.ValueOf(config))

	confObjType := confObjValue.Type()
	// getValue := reflect.ValueOf(c)

	for i := 0; i < confObjType.NumField(); i++ {
		field := confObjType.Field(i) // real key which should be get from value
		key := field.Tag.Get("configKey")
		if key == "" {
			key = field.Name
		} // use filed name instead
		configFile, ok := appValue[key]
		if ok {
			configFileValue := reflect.ValueOf(configFile)
			confField := confObjValue.Field(i)
			if confField.CanSet() {
				// switch on obj type
				kind := confField.Type().Kind()
				switch kind {
				case reflect.Slice:
					unmarshalSlice(configFileValue, confField)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
					v := processNumber(configFileValue, kind)
					confField.Set(reflect.ValueOf(v))
				default:
					confField.Set(configFileValue)
				}
			}
		}
	}
	return nil
}

// unmarshal data to obj, obj holds true data type
// while data needs type assertion
func unmarshalSlice(data reflect.Value, obj reflect.Value) {
	eleKind := obj.Type().Elem().Kind()
	switch eleKind {
	// TODO: add different types of slice
	case reflect.String:
		dataSlice, _ := data.Interface().([]interface{})
		tmpSlice := reflect.ValueOf([]string{})
		for _, value := range dataSlice {
			stringValue, ok := value.(string)
			if ok {
				tmpSlice = reflect.Append(tmpSlice, reflect.ValueOf(stringValue))
			}
		}
		obj.Set(tmpSlice)
	}
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
