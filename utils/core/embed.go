package utils

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// EmbedOptions - Embed mesaj taslağı için seçenekler
type EmbedOptions struct {
	Title         string
	Description   string
	Color         int
	Fields        []*discordgo.MessageEmbedField
	FooterText    string // Footer yazısı
	ThumbnailURL  string // Sağ üst köşedeki büyük ikon
	AuthorName    string // Author adı (üst sol)
	AuthorIconURL string // Author ikonu (üst sol) - Gözükmesi için AuthorName'de tanımlı olmak zorunda !
	ImageURL      string // Ortadaki büyük resim
	URL           string // Embed'e tıklayınca gidilecek link
}

// BuildEmbed - Embed mesajı üretir
func BuildEmbed(opts EmbedOptions, botAvatarURL string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       opts.Title,
		Description: opts.Description,
		Color:       opts.Color,
		Fields:      opts.Fields,
		URL:         opts.URL,
	}

	// Footer (en alt - sol)
	embed.Footer = &discordgo.MessageEmbedFooter{
		IconURL: botAvatarURL,
	}

	if opts.FooterText != "" {
		embed.Footer.Text = opts.FooterText
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

	checkUrlFormats(&opts)

	return s.ChannelMessageSendEmbed(channelID, BuildEmbed(opts, GetOrDefault(GetBotAvatarURL(s), Common().Icons["bot"])))
}

// SendEmbedReply, belirtilen mesajın referansını kullanarak EmbedOptions ile bir embed mesajı gönderir.
func SendEmbedReply(s *discordgo.Session, m *discordgo.MessageCreate, opts EmbedOptions) (msg *discordgo.Message, err error) {

	checkUrlFormats(&opts)

	return s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed:     BuildEmbed(opts, GetOrDefault(GetBotAvatarURL(s), Common().Icons["bot"])),
		Reference: m.Reference(),
	})
}

// Embed - Kullanıcıya özelden Embed mesajı gönderir
func SendEmbedDM(s *discordgo.Session, userID string, opts EmbedOptions) (msg *discordgo.Message, err error) {

	checkUrlFormats(&opts)

	dm, err := s.UserChannelCreate(userID)
	if err != nil {
		return nil, err
	}
	return s.ChannelMessageSendEmbed(dm.ID, BuildEmbed(opts, GetOrDefault(GetBotAvatarURL(s), Common().Icons["bot"])))
}

// Embed - Slash komutları ile belirtilen mesaja gizli (yalnızca tetikleyen görür) Embed mesajı gönderir
func SendEmbedEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, opts EmbedOptions) (err error) {

	checkUrlFormats(&opts)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{BuildEmbed(opts, GetOrDefault(GetBotAvatarURL(s), Common().Icons["bot"]))},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})

}

func checkUrlFormats(opts *EmbedOptions) {

	// URL Yanlış formatta ise kaldır , aksi takdirde mesaj hiç gönderilemez
	if opts.AuthorIconURL != "" {
		if !strings.HasPrefix(opts.AuthorIconURL, "http://") && !strings.HasPrefix(opts.AuthorIconURL, "https://") {
			opts.AuthorIconURL = ""
			Logger(WARNING, "Embed - Author URL geçersiz formatta , Mesajdan kaldırdıldı. (http:// veya https:// ile başlamalı)")
		}
	}

	// URL Yanlış formatta ise kaldır , aksi takdirde mesaj hiç gönderilemez

	if opts.ThumbnailURL != "" {
		if !strings.HasPrefix(opts.ThumbnailURL, "http://") && !strings.HasPrefix(opts.ThumbnailURL, "https://") {
			opts.ThumbnailURL = ""
			Logger(WARNING, "Embed - Thumbnail URL geçersiz formatta , Mesajdan kaldırdıldı. (http:// veya https:// ile başlamalı)")
		}
	}
}
