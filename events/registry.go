package events

import (
	"VoidBot/config"
	"VoidBot/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// SlashCommand — her slash komutun implement etmesi gereken interface
type SlashCommand interface {
	Name() string
	Description() string
	Options() []*discordgo.ApplicationCommandOption
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// PrefixCommand — her prefix komutun implement etmesi gereken interface
type PrefixCommand interface {
	Name() string // prefix olmadan: "warn", "ban"
	Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
}

// slashRegistry — slash komutları tutar
var slashRegistry = make(map[string]SlashCommand)

// prefixRegistry — prefix komutları tutar
var prefixRegistry = make(map[string]PrefixCommand)

// ── KAYIT FONKSİYONLARI ──────────────────────────────────────────────

// RegisterSlash — slash komutu kaydeder
func RegisterSlash(cmd SlashCommand) {
	slashRegistry[cmd.Name()] = cmd
	utils.Logger(utils.INFO, "Slash komutu kaydedildi: /"+cmd.Name())
}

// RegisterPrefix — prefix komutu kaydeder
func RegisterPrefix(cmd PrefixCommand) {
	prefixRegistry[cmd.Name()] = cmd
	utils.Logger(utils.INFO, "Prefix komutu kaydedildi: "+config.AppConfig.Bot.Prefix+cmd.Name())
}

// ── HANDLER FONKSİYONLARI ────────────────────────────────────────────
// Bu fonksiyonlar Bus'a kayıt olur

// HandleInteraction — gelen slash interaction'ı ilgili komuta yönlendirir
func HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	name := i.ApplicationCommandData().Name
	cmd, exists := slashRegistry[name]
	if !exists {
		return
	}

	cmd.Execute(s, i)
}

// HandleMessage — gelen mesajı prefix kontrolü yaparak ilgili komuta yönlendirir
func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Bot'un kendi mesajlarını yoksay
	if m.Author.ID == s.State.User.ID {
		return
	}

	prefix := config.AppConfig.Bot.Prefix

	// Prefix ile başlıyor mu?
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	// args oluştur: ["warn", "@kişi", "sebep"]  — prefix olmadan
	args := strings.Fields(strings.TrimPrefix(m.Content, prefix))
	if len(args) == 0 {
		return
	}

	// Komut adına göre bul ve çalıştır
	cmd, exists := prefixRegistry[args[0]]
	if !exists {
		return
	}

	cmd.Execute(s, m, args)
}

// ── SYNC ─────────────────────────────────────────────────────────────

// SyncSlashCommands — kayıtlı slash komutları Discord'a bildirir
// Bot açılışında bir kere çağrılır
func SyncSlashCommands(s *discordgo.Session) error {
	for _, cmd := range slashRegistry {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
			Name:        cmd.Name(),
			Description: cmd.Description(),
			Options:     cmd.Options(),
		})
		if err != nil {
			return err
		}
		utils.Logger(utils.OK, "Discord'a sync edildi: /"+cmd.Name())
	}
	return nil
}
