package utils

import "github.com/bwmarrin/discordgo"

//  Embed - Embed mesaj taslağı
type EmbedOptions struct {
	Title         string
	Description   string
	Color         int
	Fields        []*discordgo.MessageEmbedField
	Footer        string
	ThumbnailURL  string // Sağ üst köşedeki .png ikon URL'si
	AuthorIconURL string // Başlığın solundaki minik .png ikon URL'si
}

// Embed - Embed mesajı üretir
func BuildEmbed(opts EmbedOptions) *discordgo.MessageEmbed {

	// Temel embedi hazırla
	embed := &discordgo.MessageEmbed{
		Title:       opts.Title,
		Description: opts.Description,
		Color:       opts.Color,
		Fields:      opts.Fields,
	}

	// İmza yazısı ekle (varsa)
	if opts.Footer != "" {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: opts.Footer,
		}
	}

	// Author ikon ekle (varsa)
	if opts.AuthorIconURL != "" {
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    opts.Title,
			IconURL: opts.AuthorIconURL,
		}
		embed.Title = "" // Author varken normal title'ı çakışma olmasın diye temizle
	}

	// Thumbnail ikon ekle (varsa)
	if opts.ThumbnailURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: opts.ThumbnailURL,
		}
	}

	return embed
}

// Embed - Embed mesajı kanala gönderir
func SendEmbed(s *discordgo.Session, channelID string, opts EmbedOptions) (msg *discordgo.Message, err error) {
	return s.ChannelMessageSendEmbed(channelID, BuildEmbed(opts))
}

// Embed - Kullanıcıya özelden Embed mesajı gönderir
func SendEmbedDM(s *discordgo.Session, userID string, opts EmbedOptions) (msg *discordgo.Message, err error) {
	dm, err := s.UserChannelCreate(userID)
	if err != nil {
		return nil, err
	}

	return s.ChannelMessageSendEmbed(dm.ID, BuildEmbed(opts))
}

// Embed - Slash komutları ile belirtilen mesaja gizli (yalnızca tetikleyen görür) Embed mesajı gönderir
func SendEmbedEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, opts EmbedOptions) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{BuildEmbed(opts)},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}
