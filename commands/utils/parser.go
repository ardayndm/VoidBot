package utils

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Mention formatındaki kullanıcı ID'sini ayrıştırır.
// <@123456789012345678> -> "123456789012345678"
func ParseUserID(mention string) string {
	id := strings.TrimPrefix(mention, "<@")
	id = strings.TrimPrefix(id, "!")
	id = strings.TrimSuffix(id, ">")
	return id
}

// Kullanıcı ID'sini alır ve Discord API'si üzerinden kullanıcı bilgilerini getirir.
// Eğer kullanıcı bulunamazsa veya bir hata oluşursa, hata döner.
func GetUser(s *discordgo.Session, userID string) (*discordgo.User, error) {
	user, err := s.User(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Kullanıcı ID'sini ve sunucu ID'sini alır, ardından Discord API'si üzerinden üye bilgilerini getirir.
// Eğer üye bulunamazsa veya bir hata oluşursa, hata döner.
func GetMember(s *discordgo.Session, guildID, userID string) (*discordgo.Member, error) {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return nil, err
	}

	return member, nil
}

// Komut argümanlarından sebep metnini ayrıştırır.
// ["!warn","<@123>","spam","yapıyor"] -> "spam yapıyor"
func ParseReason(args []string, startIndex int) string {
	if len(args) <= startIndex {
		return "Sebep belirtilmedi."
	}
	return strings.Join(args[startIndex:], " ")
}

// Mesajda mention/kişi var mı kontrol eder
func HasMention(args []string) bool {
	if len(args) < 2 {
		return false
	}

	return strings.HasPrefix(args[1], "<@")
}

// Sunucu bilgilerini getirir
func GetGuild(s *discordgo.Session, guildID string) (*discordgo.Guild, error) {
	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}
	return guild, nil
}
