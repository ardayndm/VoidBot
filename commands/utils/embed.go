package utils

import "github.com/bwmarrin/discordgo"

type EmbedOptions struct {
	Title       string
	Description string
	Color       int
	Fields      []*discordgo.MessageEmbedField
	Footer      string
}

var (
	ColorPurple = 0x8B00FF
	ColorDark   = 0x2C2C2C
	ColorChain  = 0x708090
	ColorStorm  = 0x4B0082
	ColorRed    = 0xFF0000
	ColorGreen  = 0x00FF00
	ColorYellow = 0xffbb00
)

// BuildEmbed, EmbedOptions kullanarak bir MessageEmbed oluşturur.
func BuildEmbed(opts EmbedOptions) *discordgo.MessageEmbed {

	embed := &discordgo.MessageEmbed{
		Title:       opts.Title,
		Description: opts.Description,
		Color:       opts.Color,
		Fields:      opts.Fields,
	}

	if opts.Footer != "" {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: opts.Footer,
		}
	}
	return embed
}

// SendEmbed, belirtilen kanal ID'sine EmbedOptions kullanarak bir embed mesajı gönderir.
func SendEmbed(s *discordgo.Session, channelID string, opts EmbedOptions) {
	s.ChannelMessageSendEmbed(channelID, BuildEmbed(opts))
}

// SendEmbedReply, belirtilen mesajın referansını kullanarak EmbedOptions ile bir embed mesajı gönderir.
func SendEmbedReply(s *discordgo.Session, m *discordgo.MessageCreate, opts EmbedOptions) {
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed:     BuildEmbed(opts),
		Reference: m.Reference(),
	})
}

// SendEmbedDM, belirtilen kullanıcı ID'sine EmbedOptions kullanarak bir embed mesajı gönderir.
func SendEmbedDM(s *discordgo.Session, userID string, opts EmbedOptions) error {
	dm, err := s.UserChannelCreate(userID)
	if err != nil {
		return err
	}
	_, err = s.ChannelMessageSendEmbed(dm.ID, BuildEmbed(opts))
	return err
}

// SendEmbedEphemeral, belirtilen etkileşime EmbedOptions kullanarak gizli bir embed mesajı gönderir.
// Slash komutları ile kullanılmak üzere tasarlanmıştır.
func SendEmbedEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, opts EmbedOptions) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{BuildEmbed(opts)},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}
