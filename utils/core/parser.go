package utils

import "github.com/bwmarrin/discordgo"

// Parser - Sunucu bilgilerini getirir
func GetGuild(s *discordgo.Session, guildID string) (*discordgo.Guild, error) {
	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}
	return guild, nil
}

// Parser - Sunucu adını döndürür (DM'de boş string)
func GetGuildName(s *discordgo.Session, guildID string) (string, error) {
	if guildID == "" {
		return "", nil // DM
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		return "", err
	}
	return guild.Name, nil
}

// GetBotAvatarURL - Botun profil resminin URL'sini döndürür
func GetBotAvatarURL(s *discordgo.Session) string {
	user, err := s.User("@me")
	if err != nil {
		return "" // Hata olursa boş döndür
	}
	return user.AvatarURL("32") // 32x32 boyutunda (footer için ideal)
}

// getGuildIcon — target'a göre guild ikonunu çeker
func GetGuildIcon(s *discordgo.Session, t Target) string {
	if t.Interaction != nil {
		guild, err := GetGuild(s, t.Interaction.GuildID)
		if err != nil {
			return ""
		}
		return guild.IconURL("128")
	}
	if t.Message != nil {
		guild, err := GetGuild(s, t.Message.GuildID)
		if err != nil {
			return ""
		}
		return guild.IconURL("128")
	}
	return ""
}
