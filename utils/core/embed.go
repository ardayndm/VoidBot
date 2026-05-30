package utils

import "github.com/bwmarrin/discordgo"

// EmbedOptions - Embed mesaj taslağı için seçenekler
type EmbedOptions struct {
	Title         string
	Description   string
	Color         int
	Fields        []*discordgo.MessageEmbedField
	FooterText    string // Footer yazısı
	FooterIconURL string // Footer ikonu (sol alt)
	ThumbnailURL  string // Sağ üst köşedeki büyük ikon
	AuthorName    string // Author adı (üst sol)
	AuthorIconURL string // Author ikonu (üst sol)
	ImageURL      string // Ortadaki büyük resim
	URL           string // Embed'e tıklayınca gidilecek link
}

// BuildEmbed - Embed mesajı üretir
func BuildEmbed(opts EmbedOptions) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       opts.Title,
		Description: opts.Description,
		Color:       opts.Color,
		Fields:      opts.Fields,
		URL:         opts.URL,
	}

	// Footer (en alt - sol)
	if opts.FooterText != "" || opts.FooterIconURL != "" {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text:    opts.FooterText,
			IconURL: opts.FooterIconURL,
		}
	}

	// Author (en üst - sol)
	if opts.AuthorName != "" || opts.AuthorIconURL != "" {
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    opts.AuthorName,
			IconURL: opts.AuthorIconURL,
		}
	}

	// Thumbnail (sağ üst)
	if opts.ThumbnailURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: opts.ThumbnailURL,
		}
	}

	// Image (ortada - büyük resim)
	if opts.ImageURL != "" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: opts.ImageURL,
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
