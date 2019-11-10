package pro

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type ServerConfig struct {
	Auth string
	Port int32
}

var Config *ServerConfig

func LoadFile(path string) *ServerConfig {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	Config = new(ServerConfig)
	jsonError := json.Unmarshal(bytes, Config)
	if jsonError != nil {
		panic(errors.New(fmt.Sprintf("Invalid json file, please check your format. Detail is %s", jsonError.Error())))
	}
	if Config.Port <= 0 || Config.Port > 65535 {
		panic(errors.New("port should between 0 and 65535. "))
	}
	return Config
}
