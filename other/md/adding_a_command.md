# Yeni Komut Eklemek

Bu sayfada `/ban` komutunu örnek alarak yeni bir komut ekleme adımlarını göreceksin.

---

## 1. YAML dosyasını oluştur

`locale/tr/commands/ban.yaml`:

```yaml
# ==================== BAŞLIKLAR ====================
titles:
  success: "✅ Kullanıcı Yasaklandı"
  error: "❌ Yasaklama Başarısız"

# ==================== MESAJLAR ====================
# {placeholder} şeklinde dinamik değişkenler kullanabilirsin
messages:
  success: "**{user}** sunucudan yasaklandı."
  dm: "**{guild}** sunucusundan yasaklandınız."
  reason: "Sebep: {reason}"

# ==================== SEÇENEKLER ====================
# Sayı, bool, string desteklenir
options:
  cooldown: 5 # Saniye cinsinden bekleme süresi
  enabled: true # Komut aktif mi
```

> **Not:** Sadece ihtiyacın olan alanları ekle. `titles`, `messages`, `options` hepsini kullanmak zorunda değilsin.

---

## 2. Komut dosyasını oluştur

`commands/moderation/ban.go`:

```go
package moderation

import (
    events "VoidBot/events/core"
    utils "VoidBot/utils/core"
    "github.com/bwmarrin/discordgo"
)

// ── SLASH (/ban) ──────────────────────────────────────────────────────

type BanSlash struct{}

func (b BanSlash) Name() string        { return "ban" }
func (b BanSlash) Description() string { return "Bir kullanıcıyı sunucudan yasakla." }

func (b BanSlash) Options() []*discordgo.ApplicationCommandOption {
    return []*discordgo.ApplicationCommandOption{
        {
            Type:        discordgo.ApplicationCommandOptionUser,
            Name:        "kullanici",
            Description: "Yasaklanacak kullanıcı",
            Required:    true,
        },
        {
            Type:        discordgo.ApplicationCommandOptionString,
            Name:        "sebep",
            Description: "Yasaklama sebebi",
            Required:    false,
        },
    }
}

func (b BanSlash) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
    // YAML'ı yükle
    ban, err := utils.LoadCommand("ban")
    if err != nil {
        utils.RespondErrorEmbed(s, utils.Target{Interaction: i, Ephemeral: true}, "db_error", "VoidBot")
        return
    }

    // Komut aktif mi?
    if !utils.OptionBool(ban, "enabled", true) {
        return
    }

    // Parametreleri al
    opts := i.ApplicationCommandData().Options
    targetUser := opts[0].UserValue(s)

    reason := "Sebep belirtilmedi."
    if len(opts) > 1 {
        reason = opts[1].StringValue()
    }

    // Başarı mesajı
    msg := utils.FormatKeys(ban.Messages["success"], map[string]string{
        "user": targetUser.Username,
    })

    utils.RespondEmbed(s, utils.Target{Interaction: i}, utils.EmbedOptions{
        Title:       ban.Titles["success"],
        Description: msg,
        Color:       utils.ColorsInt.Success,
    })
}

// ── PREFIX (v!ban) ────────────────────────────────────────────────────

type BanPrefix struct{}

func (b BanPrefix) Name() string { return "ban" }

func (b BanPrefix) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
    // args[0] = "ban"
    // args[1] = "@kullanıcı"
    // args[2..] = sebep
}
```

---

## 3. Komutları kaydet

`main.go` içindeki `main()` fonksiyonuna ekle:

```go
events.RegisterSlash(moderation.BanSlash{})
events.RegisterPrefix(moderation.BanPrefix{})
```

---

## Kullanışlı Fonksiyonlar

### Embed Gönderme

```go
// Sadece komut kullanan görür (ephemeral)
utils.RespondErrorEmbed(s, utils.Target{
    Interaction: i,
    Ephemeral:   true,
}, "no_perm", "VoidBot")

// Kanala reply olarak gönder
utils.RespondSuccessEmbed(s, utils.Target{
    Message: m,
}, "default", "VoidBot")
```

### YAML Okuma

```go
ban, err := utils.LoadCommand("ban")

// Başlık al
ban.Titles["success"]

// Mesaj al + placeholder replace
msg := utils.FormatKeys(ban.Messages["success"], map[string]string{
    "user":   targetUser.Username,
    "reason": reason,
})

// Option oku
cooldown := utils.OptionInt(ban, "cooldown", 5)
enabled  := utils.OptionBool(ban, "enabled", true)
```

### Veritabanı

```go
ctx, cancel := utils.Ctx(5)
defer cancel()

storage.GetDB().Exec(ctx,
    "INSERT INTO bans (guild_id, user_id, reason) VALUES (?, ?, ?)",
    guildID, userID, reason,
)
```

---

## Otomatik Handler (Küfür Filtresi vb.)

Slash veya prefix olmayan, her mesajı dinleyen handler:

`commands/auto/filter.go`:

```go
package auto

import "github.com/bwmarrin/discordgo"

func Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
    // Her mesajda çalışır
    // Küfür kontrolü, XP verme vb.
}
```

`main.go`'da:

```go
events.Bus.OnMessage(auto.Handle)
```

---

## Örnek YAML Şablonu

Yeni komut eklerken `locale/example/commands/command_example.yaml` dosyasını şablon olarak kullanabilirsin.
