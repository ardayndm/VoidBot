# Lokalizasyon

VoidBot, tüm metinleri YAML dosyalarında tutar. Kod değiştirmeden dil eklenebilir, metinler düzenlenebilir.

---

## Klasör Yapısı

```
locale/
└── tr/
    ├── common.yaml          # Genel mesajlar (tüm komutlarda kullanılır)
    └── commands/
        ├── warn.yaml        # Sadece warn komutuna özel
        ├── ban.yaml         # Sadece ban komutuna özel
        └── ...
```

---

## common.yaml

Tüm komutlarda ortak kullanılan mesajlar burada durur.

```yaml
error:
  title: "❌ Bir Hata Oluştu"
  messages:
    no_perm: "Bu komutu kullanmak için gerekli yetkiye sahip değilsiniz."
    user_not_found: "Belirtilen kullanıcı sunucuda bulunamadı."

warning:
  title: "⚠️ Uyarı"
  messages:
    cooldown: "Bu komutu tekrar kullanmak için **{seconds} saniye** beklemelisiniz."

success:
  title: "✅ Başarılı"

info:
  title: "ℹ️ Bilgi"

icons:
  error: "https://..." # URL olmak zorunda, emoji geçersiz
  success: "https://..."
```

> **Not:** `icons` altındaki değerler mutlaka `http://` veya `https://` ile başlamalıdır. Emoji veya boş string Discord tarafından reddedilir.

---

## Komut YAML'ı

Her komut kendi mesajlarını kendi dosyasında tutar.

```yaml
# locale/tr/commands/warn.yaml

# ==================== BAŞLIKLAR ====================
titles:
  success: "✅ Kullanıcı Uyarıldı"
  error: "❌ Uyarma Başarısız"

# ==================== MESAJLAR ====================
# {placeholder} şeklinde dinamik değişkenler kullanabilirsin
messages:
  success: "**{user}** uyarıldı."
  dm: "**{guild}** sunucusundan uyarı aldınız."
  reason: "Sebep: {reason}"

# ==================== SEÇENEKLER ====================
# Sayı, bool, string desteklenir
options:
  cooldown: 5 # Saniye cinsinden bekleme süresi
  enabled: true # Komut aktif mi
```

**Kural:**

- `common.yaml` → birden fazla komutta çıkabilecek mesajlar
- `commands/warn.yaml` → sadece warn komutuna özel mesajlar

---

## Kod Tarafında Kullanım

```go
// Komut YAML'ını yükle (ilk seferinde diskten, sonra cache'den)
warn, err := utils.LoadCommand("warn")

// Mesaj al
warn.Messages["success"]

// Option oku
cooldown  := utils.OptionInt(warn, "cooldown", 5)
enabled   := utils.OptionBool(warn, "enabled", true)

// Ortak mesajlar
utils.Common().Error.Messages["no_perm"]
utils.Common().Icons["success"]
```

---

## {Placeholder} Sistemi

```go
text := utils.FormatKeys("**{guild}** sunucusunda bir uyarı aldın.", map[string]string{
    "guild": "Void",
})
// → "**Void** sunucusunda bir uyarı aldın."
```

---

## Yeni Dil Eklemek

```
locale/
├── tr/
│   ├── common.yaml
│   └── commands/warn.yaml
└── en/              ← yeni dil klasörü
    ├── common.yaml
    └── commands/warn.yaml
```

`config/core/bot.yaml` dosyasında dili değiştir:

```yaml
lang: "en"
```
