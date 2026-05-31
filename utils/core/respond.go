package utils

import (
	"github.com/bwmarrin/discordgo"
)

// Respond — mesajın nereye gideceğini tutar
type Target struct {
	// İkisinden biri dolu olur, diğeri nil
	Message     *discordgo.MessageCreate
	Interaction *discordgo.InteractionCreate
	Ephemeral   bool
}

// Respond - Target'tan guild ID'sini alır
func (t *Target) GetGuildID() string {
	if t.Interaction != nil {
		return t.Interaction.GuildID
	}
	if t.Message != nil {
		return t.Message.GuildID
	}
	return ""
}

// Respond - Mesaj veya interaction'a embed gönderir (DM İçermez !)
func RespondEmbed(s *discordgo.Session, t Target, opts EmbedOptions) error {

	// Mesaj tipi Etkileşimli/slash mesaj ise
	if t.Interaction != nil {
		return SendEmbedEphemeral(s, t.Interaction, opts)
	}

	// Mesaj tipi normal mesaj ise
	if t.Message != nil {
		_, err := SendEmbedReply(s, t.Message, opts)
		return err
	}

	return nil
}
