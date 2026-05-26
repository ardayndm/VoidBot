package utils

import (
	"voidbot/config"

	"github.com/bwmarrin/discordgo"
)

// Discord'da her işlemin bir permission bit'i vardır.

func HasPermission(s *discordgo.Session, guildID, userID string, permission int64) (bool, error) {
	// Önce sunucu sahibi mi bakalım, sahibin her yetkisi vardır

	guild, err := s.Guild(guildID)

	if err != nil {
		return false, err
	}

	if guild.OwnerID == userID {
		return true, nil
	}

	// Member'in rollerini çek
	member, err := GetMember(s, guildID, userID)
	if err != nil {
		return false, err
	}

	// Her rolün permissions'ını kontrol et
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			continue
		}

		// Bitwise AND - o yetki bu rolde var mı?
		if role.Permissions&permission != 0 {
			return true, nil
		}
	}
	return false, nil
}

func IsAdmin(s *discordgo.Session, guildID, userID string) (bool, error) {
	return HasPermission(s, guildID, userID, discordgo.PermissionAdministrator)
}

func CanKick(s *discordgo.Session, guildID, userID string) (bool, error) {
	return HasPermission(s, guildID, userID, discordgo.PermissionKickMembers)
}

func CanBan(s *discordgo.Session, guildID, userID string) (bool, error) {
	return HasPermission(s, guildID, userID, discordgo.PermissionBanMembers)
}

func CanManageMessages(s *discordgo.Session, guildID, userID string) (bool, error) {
	return HasPermission(s, guildID, userID, discordgo.PermissionManageMessages)
}

func SendNoPermissionEmbedMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	SendEmbedReply(s, m, EmbedOptions{
		Title:       "Yetersiz İzin",
		Description: "Bu komutu kullanmak için yeterli izniniz yok.",
		Color:       ColorRed, // Kırmızı renk
		Footer:      config.Bot.Name,
	})
}
