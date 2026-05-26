package utils

import (
	"voidbot/config"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrNoPermission = "Bu komutu kullanmak için yetkin yok."
	ErrNoMention    = "Lütfen bir kullanıcı etiketle."
	ErrUserNotFound = "Kullanıcı bulunamadı."
	ErrMissingArgs  = "Eksik argüman. Kullanım: "
	ErrSelfAction   = "Kendine bu işlemi yapamazsın."
	ErrBotAction    = "Bot hesaplarına bu işlemi yapamazsın."
	ErrHigherRole   = "Bu kullanıcının rolü senden yüksek."
)

// Hata mesajını embed olarak gönderir
func SendError(s *discordgo.Session, m *discordgo.MessageCreate, description string) {
	SendEmbedReply(s, m, EmbedOptions{
		Title:       "Hata",
		Description: description,
		Color:       ColorRed,
		Footer:      config.Bot.Name,
	})
}

func SendSuccess(s *discordgo.Session, m *discordgo.MessageCreate, description string) {
	SendEmbedReply(s, m, EmbedOptions{
		Title:       "Başarılı",
		Description: description,
		Color:       ColorGreen,
		Footer:      config.Bot.Name,
	})
}

func SendWarning(s *discordgo.Session, m *discordgo.MessageCreate, description string) {
	SendEmbedReply(s, m, EmbedOptions{
		Title:       "Uyarı",
		Description: description,
		Color:       ColorYellow,
		Footer:      config.Bot.Name,
	})
}
