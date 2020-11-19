package tool

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _database *gorm.DB

func GetDatabase() *gorm.DB {
	if _database == nil {
		InitDatabase()
	}
	return _database
}

func InitDatabase() {
	config := GetConfig()

	dsn := config.DBUsername + ":" + config.DBPassword +
		"@tcp(" + config.DBAddr + ")/zhihu?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	_database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err.Error())
	}

}
