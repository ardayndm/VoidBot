package utils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// embedType - Embed tipi (error, success, warning, info)
type embedType string

const (
	embedError   embedType = "error"
	embedSuccess embedType = "success"
	embedWarning embedType = "warning"
	embedInfo    embedType = "info"
)

// respondEmbedByType - Genel embed gönderme fonksiyonu
func respondEmbedByType(s *discordgo.Session, t Target, messageKey, botName string, eType embedType) error {
	// Common verilerini al
	common := Common()

	// Embed tipine göre verileri seç
	var title, message, iconKey string
	var color int

	switch eType {
	case embedError:
		title = common.Error.Title
		message = getMessage(common.Error.Messages, messageKey)
		iconKey = "error"
		color = ColorsInt.Error
	case embedSuccess:
		title = common.Success.Title
		message = getMessage(common.Success.Messages, messageKey)
		iconKey = "success"
		color = ColorsInt.Success
	case embedWarning:
		title = common.Warning.Title
		message = getMessage(common.Warning.Messages, messageKey)
		iconKey = "warning"
		color = ColorsInt.Warning
	case embedInfo:
		title = common.Info.Title
		message = getMessage(common.Info.Messages, messageKey)
		iconKey = "info"
		color = ColorsInt.Info
	default:
		return fmt.Errorf("bilinmeyen embed tipi: %s", eType)
	}

	if message == "" {
		message = "ERR"
	}

	// Embed seçeneklerini hazırla
	embedOpt := EmbedOptions{
		Title:        title,
		Description:  message,
		Color:        color,
		ThumbnailURL: common.Icons[iconKey],
		FooterText:   GetOrDefault(botName, "VoidBot"),
	}

	// Sunucu adını al (varsa)
	if guildName, err := GetGuildName(s, t.GetGuildID()); err == nil && guildName != "" {
		embedOpt.AuthorName = guildName
	}

	// İkonları ayarla (sunucu ikonu yoksa bot ikonu kullan)
	embedOpt.AuthorIconURL = GetOrDefault(GetGuildIcon(s, t), common.Icons["bot"])
	return RespondEmbed(s, t, embedOpt)
}

// getMessage - Mesaj map'inden mesajı alır, yoksa hata loglar
func getMessage(messages map[string]string, key string) string {
	if msg, exists := messages[key]; exists && msg != "" {
		return msg
	}
	Logger(ERROR, fmt.Sprintf("Mesaj anahtarı bulunamadı: %s", key))
	return ""
}

// getOrDefault - Değer boşsa varsayılanı döndürür
func GetOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// ==================== PUBLIC FONKSİYONLAR ====================

// RespondErrorEmbed — hata embed'i gönderir
func RespondErrorEmbed(s *discordgo.Session, t Target, messageKey, botName string) error {
	return respondEmbedByType(s, t, messageKey, botName, embedError)
}

// RespondSuccessEmbed — başarı embed'i gönderir
func RespondSuccessEmbed(s *discordgo.Session, t Target, messageKey, botName string) error {
	return respondEmbedByType(s, t, messageKey, botName, embedSuccess)
}

// RespondWarningEmbed — uyarı embed'i gönderir
func RespondWarningEmbed(s *discordgo.Session, t Target, messageKey, botName string) error {
	return respondEmbedByType(s, t, messageKey, botName, embedWarning)
}

// RespondInfoEmbed — bilgi embed'i gönderir
func RespondInfoEmbed(s *discordgo.Session, t Target, messageKey, botName string) error {
	return respondEmbedByType(s, t, messageKey, botName, embedInfo)
}
