package tool

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

var Logger *logrus.Logger
var Formatter logrus.Formatter

func InitLogger() {
	if Logger != nil {
		return
	}

	Logger = logrus.New()
	Formatter = &logrus.TextFormatter{
		TimestampFormat:"2006-01-02 15:04:05",
		DisableQuote: true,
	}
	Logger.SetFormatter(Formatter)

	if Cfg.Mode == "debug" {
		Logger.SetLevel(logrus.DebugLevel)
	}

	if Cfg.LogFile != "" {
		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err!= nil {
			fmt.Println("Error: ", err.Error())
		}
		Logger.Out = src


		apiLogPath := path.Join(Cfg.LogPath, Cfg.LogFile)
		logWriter, err := rotatelogs.New(
			apiLogPath + ".%Y-%m-%d-%H-%M.log",
			rotatelogs.WithLinkName(Cfg.LogFile),      // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
			rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
		)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}
		writeMap := lfshook.WriterMap{
			logrus.InfoLevel:  logWriter,
			logrus.FatalLevel: logWriter,
		}
		lfHook := lfshook.NewHook(writeMap, Formatter)
		Logger.AddHook(lfHook)
	}
}
