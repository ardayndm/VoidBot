package utils

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Respond — mesajın nereye gideceğini tutar
type Target struct {
	// İkisinden biri dolu olur, diğeri nil
	Message     *discordgo.MessageCreate
	Interaction *discordgo.InteractionCreate
	Ephemeral   bool
}

// Respond - Mesaj veya interaction'a embed gönderir
func RespondEmbed(s *discordgo.Session, t Target, opts EmbedOptions) error {

	// URL Yanlış formatta ise kaldır , aksi takdirde mesaj hiç gönderilemez
	authorUrl := opts.AuthorIconURL
	if !strings.HasPrefix(authorUrl, "http://") && !strings.HasPrefix(authorUrl, "https://") {
		opts.AuthorIconURL = ""
		Logger(WARNING, "Respond - Author URL geçersiz formatta , Mesajdan kaldırdıldı. (http:// veya https:// ile başlamalı)")
	}

	// URL Yanlış formatta ise kaldır , aksi takdirde mesaj hiç gönderilemez
	thumbUrl := opts.ThumbnailURL
	if !strings.HasPrefix(thumbUrl, "http://") && !strings.HasPrefix(thumbUrl, "https://") {
		opts.AuthorIconURL = ""
		Logger(WARNING, "Respond - Thumbnail URL geçersiz formatta , Mesajdan kaldırdıldı. (http:// veya https:// ile başlamalı)")
	}

	// Mesaj tipi Etkileşimli/slash mesaj ise
	if t.Interaction != nil {
		flags := discordgo.MessageFlags(0)
		if t.Ephemeral {
			flags = discordgo.MessageFlagsEphemeral
		}
		err := s.InteractionRespond(t.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{BuildEmbed(opts)},
				Flags:  flags,
			},
		})
		return err
	}

	// Mesaj tipi normal mesaj ise
	if t.Message != nil {
		_, err := s.ChannelMessageSendComplex(t.Message.ChannelID, &discordgo.MessageSend{
			Embed:     BuildEmbed(opts),
			Reference: t.Message.Reference(),
		})
		return err
	}

	return nil
}
