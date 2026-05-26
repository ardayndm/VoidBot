package utils

import (
	"log"
	"time"
)

type LogLevel string

const (
	INFO    LogLevel = "[INFO]"
	WARNING LogLevel = "[WARNING]"
	ERROR   LogLevel = "[ERROR]"
	SUCCESS LogLevel = "[OK]"
)

func logger(level LogLevel, message string) {
	logMessage := string(level) + " " + message
	log.Println(logMessage)
}

func LogInfo(message string) {
	logger(INFO, message)
}

func LogWarning(message string) {
	logger(WARNING, message)
}

func LogError(message string) {
	logger(ERROR, message)
}

func LogSuccess(message string) {
	logger(SUCCESS, message)
}

// Komut logları için özel format
// [15:04:05] [CMD] kullanıcı#1234 → !warn sebep
func LogCommand(username, command string) {
	timestamp := time.Now().Format("15:04:05")
	log.Printf("[%s] [CMD] %s → %s\n", timestamp, username, command)
}
