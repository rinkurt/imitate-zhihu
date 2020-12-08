package tool

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err!= nil{
		fmt.Println("err", err)
	}
	Logger.Out = src
	Logger.SetLevel(logrus.DebugLevel)
	formatter := &logrus.TextFormatter{
		TimestampFormat:"2006-01-02 15:04:05",
	}
	Logger.SetFormatter(formatter)
	apiLogPath := "./log/log.txt"
	logWriter, err := rotatelogs.New(
		apiLogPath+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(apiLogPath), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour), // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, formatter)
	Logger.AddHook(lfHook)
}

// 日志同时输出到终端和文件
func LogOutError(msg string) {
	Logger.Error(msg)
	fmt.Println(msg)
}

func LogOutErrorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
	fmt.Printf(format, args...)
}

func LogOutInfof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
	fmt.Printf(format, args...)
}