package moderation

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"voidbot/commands/utils"
	"voidbot/config"
	"voidbot/database"
)

type WarnCommand struct{}

func (w WarnCommand) Name() string { return "!warn" }

func (w WarnCommand) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

	// args[0] = "!warn"
	// args[1] = "@mention"
	// args[2..] = sebep kelimeleri

	// ── 1. COOLDOWN ──────────────────────────────────────────────────
	if !utils.Cooldown.Check(m.Author.ID, "warn", 5*time.Second) {
		remaining := utils.Cooldown.Remaining(m.Author.ID, "warn", 5*time.Second)
		utils.SendWarning(s, m, fmt.Sprintf("Bu komutu tekrar kullanmak için **%.0f saniye** bekle.", remaining.Seconds()))
		return
	}

	// ── 2. YETKİ ─────────────────────────────────────────────────────
	ok, err := utils.CanKick(s, m.GuildID, m.Author.ID)
	if err != nil || !ok {
		utils.SendNoPermissionEmbedMessage(s, m)
		return
	}

	// ── 3. MENTION VAR MI? ────────────────────────────────────────────
	if len(args) < 2 || !utils.HasMention(args) {
		utils.SendError(s, m, utils.ErrNoMention)
		return
	}

	// ── 4. HEDEF KULLANICI ────────────────────────────────────────────
	targetID := utils.ParseUserID(args[1])
	targetUser, err := utils.GetUser(s, targetID)
	if err != nil {
		utils.SendError(s, m, utils.ErrUserNotFound)
		return
	}

	// ── 5. KENDİNE / BOTA UYARI VERME ────────────────────────────────
	// if targetID == m.Author.ID {
	// 	utils.SendError(s, m, utils.ErrSelfAction)
	// 	return
	// }
	if targetUser.Bot {
		utils.SendError(s, m, utils.ErrBotAction)
		return
	}

	// ── 6. SEBEP ──────────────────────────────────────────────────────
	reason := utils.ParseReason(args, 2)

	// ── 7. MYSQL'E KAYDET ─────────────────────────────────────────────
	_, err = database.DB.Exec(
		`INSERT INTO warnings (guild_id, user_id, reason, moderator, created_at)
		 VALUES (?, ?, ?, ?, NOW())`,
		m.GuildID, targetID, reason, m.Author.ID,
	)
	if err != nil {
		utils.LogError("warn komutu DB hatası: " + err.Error())
		utils.SendError(s, m, "Veritabanı hatası, uyarı kaydedilemedi.")
		return
	}

	// ── 8. TOPLAM UYARI SAYISI ────────────────────────────────────────
	var totalWarnings int
	row := database.DB.QueryRow(
		`SELECT COUNT(*) FROM warnings WHERE guild_id = ? AND user_id = ?`,
		m.GuildID, targetID,
	)
	if err := row.Scan(&totalWarnings); err != nil {
		totalWarnings = 1
	}

	// ── 9. KANALA BİLDİR ──────────────────────────────────────────────
	utils.SendEmbedReply(s, m, utils.EmbedOptions{
		Title: "⚠️ Uyarı Verildi",
		Color: utils.ColorPurple,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Kullanıcı", Value: targetUser.Mention(), Inline: true},
			{Name: "Moderatör", Value: m.Author.Mention(), Inline: true},
			{Name: "Sebep", Value: reason, Inline: false},
			{Name: "Toplam Uyarı", Value: fmt.Sprintf("%d", totalWarnings), Inline: true},
		},
		Footer: config.Bot.Name,
	})

	// ── 10. DM AT ─────────────────────────────────────────────────────

	guild, err := utils.GetGuild(s, m.GuildID)
	var guildName string
	var guildIcon string // Sunucu ikon linki

	if err != nil {
		utils.LogError("Sunucu bilgisi alınamadı: " + err.Error())
		guildName = m.GuildID
	} else {
		guildName = guild.Name
		// EĞER sunucunun bir ikonu varsa URL'ini alıyoruz, yoksa boş kalıyor
		if guild.Icon != "" {
			guildIcon = guild.IconURL("128") // 128x128 boyutunda ikon URL'si
		}
	}

	if err != nil {
		utils.LogError("Sunucu bilgisi alınamadı: " + err.Error())
		guildName = m.GuildID
	} else {
		guildName = guild.Name
	}

	dmErr := utils.SendEmbedDM(s, targetID, utils.EmbedOptions{
		Title:       "⚠️ Uyarı Aldın",
		Description: fmt.Sprintf("**%s** sunucusunda bir uyarı aldın.", guildName),
		Color:       utils.ColorStorm,
		Thumbnail:   utils.WarningIconURL,
		AuthorIcon:  guildIcon,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Sebep", Value: reason, Inline: false},
			{Name: "Moderatör", Value: m.Author.Username, Inline: true},
			{Name: "Toplam Uyarın", Value: fmt.Sprintf("%d", totalWarnings), Inline: true},
		},
		Footer: config.Bot.Name,
	})

	if dmErr != nil {
		utils.LogWarning("DM gönderilemedi: " + targetUser.Username + " sebep: " + dmErr.Error())
	}

	// ── 11. 3 UYARIDA OTOMATİK BAN ───────────────────────────────────
	if totalWarnings >= 3 {
		banReason := fmt.Sprintf("Otomatik ban: %d uyarıya ulaşıldı.", totalWarnings)

		banErr := s.GuildBanCreateWithReason(m.GuildID, targetID, banReason, 0)
		if banErr != nil {
			utils.LogError("Otomatik ban başarısız " + "[GuildID : " + m.GuildID + " TargetID : " + targetID + "]: " + banErr.Error())
		} else {
			utils.SendEmbedReply(s, m, utils.EmbedOptions{
				Title:       "🔥 Otomatik Ban",
				Description: fmt.Sprintf("%s **%d uyarıya** ulaştığı için otomatik olarak banlandı.", targetUser.Mention(), totalWarnings),
				Color:       utils.ColorRed,
				Footer:      config.Bot.Name,
			})
			cleanupWarnings(m.GuildID, targetID)
		}
	}

	// ── 12. COOLDOWN SET ──────────────────────────────────────────────
	utils.Cooldown.Set(m.Author.ID, "warn")

	// ── 13. LOG ───────────────────────────────────────────────────────
	utils.LogCommand(m.Author.Username, fmt.Sprintf("warn → %s (%s)", targetUser.Username, reason))
}

func cleanupWarnings(guildID, userID string) {
	_, err := database.DB.Exec(
		`DELETE FROM warnings WHERE guild_id = ? AND user_id = ?`,
		guildID, userID,
	)
	if err != nil {
		utils.LogError("Uyarı temizleme hatası: " + err.Error())
	}
}
