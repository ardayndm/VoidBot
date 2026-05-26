package handlers

import (
	"strings"
	"voidbot/commands"
	"voidbot/config"

	"github.com/bwmarrin/discordgo"
)

func MessageReceived(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Mesajın bu bot tarafından gönderilip gönderilmediğini kontrol et
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Mesajın belirli bir önekle başlayıp başlamadığını kontrol et
	if !strings.HasPrefix(m.Content, config.Bot.Prefix) {
		return
	}

	// İçerik boş ise sessiz kal
	args := strings.Fields(m.Content)
	if len(args) == 0 {
		return
	}

	// Gelen komutu registry servisine yönlendir
	commands.Handle(s, m, args)
}
