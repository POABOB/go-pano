package utils

import (
	"path"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var (
	logPath = "./log"
	logFile = "gin.log"
)

var LogInstance = logrus.New()

// 日誌初始化
func init() {
	// 打開文件
	logFileName := path.Join(logPath, logFile)
	// 使用滾動壓縮方式記錄日誌
	rolling(logFileName)
	// 設置日誌輸出JSON格式
	LogInstance.SetFormatter(&logrus.JSONFormatter{})
	// 設置日誌記錄級別
	// LogInstance.SetLevel(logrus.DebugLevel)
	LogInstance.SetLevel(logrus.InfoLevel)
	// TODO 新增 Method URI IP
}

// 日誌滾動設置
func rolling(logFile string) {
	// 設置輸出
	LogInstance.SetOutput(&lumberjack.Logger{
		Filename:   logFile, //日誌文件位置
		MaxSize:    1,       // 單文件最大容量,單位是MB
		MaxBackups: 3,       // 最大保留過期文件個數
		MaxAge:     1,       // 保留過期文件的最大時間間隔,單位是天
		Compress:   true,    // 是否需要壓縮滾動日誌, 使用的 gzip 壓縮
	})
}
