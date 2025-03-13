package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	// 日志级别常量
	LogLevelDebug = "DEBUG"
	LogLevelInfo  = "INFO"
	LogLevelWarn  = "WARN"
	LogLevelError = "ERROR"

	// 日志文件
	logFile *os.File
	Logger  *log.Logger
)

// InitLogger 初始化日志系统
func InitLogger() error {
	// 创建日志目录
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 创建日志文件，使用当前日期作为文件名
	currentDate := time.Now().Format("2006-01-02")
	logFilePath := filepath.Join(logDir, fmt.Sprintf("multi-cloud-landing-zone-%s.log", currentDate))
	
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %w", err)
	}

	// 初始化日志记录器
	Logger = log.New(logFile, "", log.LstdFlags)
	
	// 记录初始化成功信息
	LogInfo("日志系统初始化成功，日志文件: " + logFilePath)
	
	return nil
}

// LogDebug 记录调试级别日志
func LogDebug(message string) {
	logWithLevel(LogLevelDebug, message)
}

// LogInfo 记录信息级别日志
func LogInfo(message string) {
	logWithLevel(LogLevelInfo, message)
}

// LogWarn 记录警告级别日志
func LogWarn(message string) {
	logWithLevel(LogLevelWarn, message)
}

// LogError 记录错误级别日志
func LogError(message string) {
	logWithLevel(LogLevelError, message)
}

// logWithLevel 使用指定级别记录日志
func logWithLevel(level, message string) {
	if Logger != nil {
		Logger.Printf("[%s] %s", level, message)
	} else {
		// 如果日志系统未初始化，则输出到标准输出
		log.Printf("[%s] %s", level, message)
	}
}

// CloseLogger 关闭日志文件
func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}
