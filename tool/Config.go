package tool

import (
	"bufio"
	"encoding/json"
	"os"
)

// 添加配置方法：
// 1. 在 config/xxx.json 添加项
// 2. 在此处 Config 结构添加项

type Config struct {
	Port       string `json:"port"`
	Mode       string `json:"mode"` // debug or release
	DBAddr     string `json:"db_addr"`
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
	LogFile    string `json:"log_file"`
	JwtSecret  string `json:"jwt_secret"`
}

var config *Config = nil

func GetConfig() *Config {
	if config == nil {
		var err error
		if os.Getenv("IZ_ENV_MODE") == "release" {
			err = ParseConfig("./config/release.json")
		} else {
			err = ParseConfig("./config/debug.json")
		}
		if err != nil {
			panic(err.Error())
		}
	}
	return config
}

func ParseConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	config = new(Config)
	if err := decoder.Decode(config); err != nil {
		config = nil
		return err
	}

	if config.LogFile == "" {
		config.LogFile = os.Getenv("IZ_LOG_FILE")
	}
	if config.DBAddr == "" {
		config.DBAddr = os.Getenv("IZ_DB_ADDR")
	}
	if config.DBUsername == "" {
		config.DBUsername = os.Getenv("IZ_DB_USERNAME")
	}
	if config.DBPassword == "" {
		config.DBPassword = os.Getenv("IZ_DB_PASSWORD")
	}
	if config.JwtSecret == "" {
		config.JwtSecret = os.Getenv("IZ_JWT_SECRET")
	}
	MySecret = []byte(config.JwtSecret)
	return nil
}
