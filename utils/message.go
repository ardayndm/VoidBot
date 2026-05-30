package utils

import (
	"github.com/bwmarrin/discordgo"
)

// RespondError — hata embed'i gönderir
func RespondErrorEmbed(s *discordgo.Session, t Target, messageKey string) error {
	guild := getGuildIcon(s, t)
	return RespondEmbed(s, t, EmbedOptions{
		Title:         Common().Error.Title,
		Description:   Common().Error.Messages[messageKey],
		Color:         ColorsInt.Error,
		AuthorIconURL: guild,
		ThumbnailURL:  Common().Icons["error"],
	})
}

// RespondSuccess — başarı embed'i gönderir
func RespondSuccessEmbed(s *discordgo.Session, t Target, messageKey string) error {
	guild := getGuildIcon(s, t)
	return RespondEmbed(s, t, EmbedOptions{
		Title:         Common().Success.Title,
		Description:   Common().Success.Messages[messageKey],
		Color:         ColorsInt.Success,
		AuthorIconURL: guild,
		ThumbnailURL:  Common().Icons["success"],
	})
}

// RespondWarning — uyarı embed'i gönderir
func RespondWarningEmbed(s *discordgo.Session, t Target, messageKey string) error {
	guild := getGuildIcon(s, t)
	return RespondEmbed(s, t, EmbedOptions{
		Title:         Common().Warning.Title,
		Description:   Common().Warning.Messages[messageKey],
		Color:         ColorsInt.Warning,
		AuthorIconURL: guild,
		ThumbnailURL:  Common().Icons["warning"],
	})
}

// RespondInfo — bilgi embed'i gönderir
func RespondInfoEmbed(s *discordgo.Session, t Target, messageKey string) error {
	guild := getGuildIcon(s, t)
	return RespondEmbed(s, t, EmbedOptions{
		Title:         Common().Info.Title,
		Description:   Common().Info.Messages[messageKey],
		Color:         ColorsInt.Info,
		AuthorIconURL: guild,
		ThumbnailURL:  Common().Icons["info"],
	})
}

// getGuildIcon — target'a göre guild ikonunu çeker
func getGuildIcon(s *discordgo.Session, t Target) string {
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
