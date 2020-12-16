package tool

import (
	"bufio"
	"encoding/json"
	"os"
)

// 添加配置方法：
// 1. 在 Cfg/xxx.json 添加项
// 2. 在此处 Config 结构添加项

type Config struct {
	Port          string `json:"port"`
	Mode          string `json:"mode"` // debug or release
	DBAddr        string `json:"db_addr"`
	DBUsername    string `json:"db_username"`
	DBPassword    string `json:"db_password"`
	LogPath       string `json:"log_path"`
	LogFile       string `json:"log_file"`
	JwtSecret     string `json:"jwt_secret"`
	RedisAddr     string `json:"redis_addr"`
	RedisPassword string `json:"redis_password"`
}

var Cfg *Config

func InitConfig() {
	if Cfg != nil {
		return
	}
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

func GetEnvIfExist(dest *string, key string) {
	if v := os.Getenv(key); v != "" {
		*dest = v
	}
}

func ParseConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	Cfg = new(Config)
	if err := decoder.Decode(Cfg); err != nil {
		Cfg = nil
		return err
	}

	GetEnvIfExist(&Cfg.LogFile, "IZ_LOG_FILE")
	GetEnvIfExist(&Cfg.DBAddr, "IZ_DB_ADDR")
	GetEnvIfExist(&Cfg.DBUsername, "IZ_DB_USERNAME")
	GetEnvIfExist(&Cfg.DBPassword, "IZ_DB_PASSWORD")
	GetEnvIfExist(&Cfg.JwtSecret, "IZ_JWT_SECRET")
	GetEnvIfExist(&Cfg.RedisAddr, "IZ_REDIS_ADDR")
	GetEnvIfExist(&Cfg.RedisPassword, "IZ_REDIS_PASSWORD")

	MySecret = []byte(Cfg.JwtSecret)
	return nil
}
