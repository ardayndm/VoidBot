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
