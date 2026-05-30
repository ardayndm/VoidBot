package utils

import (
	"fmt"
	"log"
	"strings"
)

type LogLevel string

// Logger - Uygulamanın aktif log seviyesi
var AppLogLevel LogLevel = none

// Logger - LogLevel'leri
const (
	none LogLevel = "NONE"
	all  LogLevel = "ALL"
)

// Logger - LogLevel'leri
const (
	INFO    LogLevel = "INFO"
	ERROR   LogLevel = "ERROR"
	WARNING LogLevel = "WARNING"
	OK      LogLevel = "OK"
)

// Logger - ANSI Renk kodları
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

// Logger - Log seviyesini ayarlar
func SetLogLevel(level string) {

	switch strings.ToUpper(level) {
	case "ALL":
		AppLogLevel = all
	case "INFO":
		AppLogLevel = INFO
	case "ERROR":
		AppLogLevel = ERROR
	case "WARNING":
		AppLogLevel = WARNING
	case "OK":
		AppLogLevel = OK
	case "NONE":
		AppLogLevel = none
	default:
		AppLogLevel = INFO // Varsayılan
	}
}

// Logger - Level'e göre renk seç
func getColor(level LogLevel) string {
	switch level {
	case INFO:
		return colorCyan
	case ERROR:
		return colorRed
	case WARNING:
		return colorYellow
	case OK:
		return colorGreen
	default:
		return colorWhite
	}
}

// Logger - Log yazdırma seviyesini kontrol eder
func isLogAvailable(level LogLevel) bool {
	if AppLogLevel == none {
		return false
	}

	if AppLogLevel == all {
		return true
	}

	if level == AppLogLevel {
		return true
	}

	return false
}

// Logger - Renkli log mesajı yazar
func Logger(level LogLevel, message string) {

	if !isLogAvailable(level) {
		return
	}

	color := getColor(level)

	// Renkli log mesajı
	logMessage := fmt.Sprintf("%s%s%s %s",
		color,
		fmt.Sprintf("[%s]", level),
		colorReset,
		message)

	log.Println(logMessage)
}
