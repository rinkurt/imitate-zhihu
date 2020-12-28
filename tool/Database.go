package tool

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func GetDatabase() *gorm.DB {
	return db
}

func InitDatabase(database string) {
	if db != nil {
		return
	}

	dsn := Cfg.DBUsername + ":" + Cfg.DBPassword +
		"@tcp(" + Cfg.DBAddr + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err.Error())
	}

}
