package tool

import (
	"bufio"
	"encoding/json"
	"os"
)

type Config struct {
	Port string `json:"port"`
	Mode string `json:"mode"` // debug or release
	DBAddr string `json:"db_addr"`
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
}

var _config *Config = nil

func GetConfig() *Config {
	if _config == nil {
		_, err := ParseConfig("./config/app.json")
		if err != nil {
			panic(err.Error())
		}
	}
	return _config
}

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	_config = new(Config)
	if err := decoder.Decode(_config); err != nil {
		_config = nil
		return nil, err
	}
	return _config, nil
}