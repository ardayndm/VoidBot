# Proje Yapısı

```
VoidBot/
│
├── main.go                        # Giriş noktası
│
├── config/core/
│   ├── config.go                  # .env okuma, yapılandırma yükleme
│   ├── bot.yaml                   # Bot adı, prefix, dil
│   └── color.yaml                 # Embed renkleri
│
├── cache/
│   ├── interface.go               # Cache arayüzü
│   └── redis.go                   # Redis implementasyonu
│
├── database/
│   ├── core/interface.go          # Database arayüzü
│   └── mysql.go                   # MySQL implementasyonu
│
├── storage/core/
│   ├── storage.go                 # DB + Cache tek çatı altında
│   └── migrator.go                # SQL migration yöneticisi
│
├── events/core/
│   ├── bus.go                     # Event dağıtıcı
│   └── registry.go                # Komut kayıt ve yönlendirme sistemi
│
├── commands/
│   ├── moderation/
│       ├── example.go             # komutlar
│       └── ...
│
├── utils/core/
│   ├── color.go                   # Renk yükleme
│   ├── context.go                 # Timeout context
│   ├── embed.go                   # Embed oluşturucu
│   ├── format.go                  # {placeholder} sistemi
│   ├── localization.go            # YAML dil sistemi
│   ├── logger.go                  # Renkli terminal logu
│   ├── message.go                 # Hazır embed gönderici
│   ├── parser.go                  # Discord API yardımcıları
│   ├── read.go                    # YAML okuyucu
│   └── respond.go                 # Mesaj/Interaction gönderici
│
├── locale/
│   └── tr/
│       ├── common.yaml            # Genel mesajlar
│       └── commands/
│           ├── example.yaml
|           └── ...
│
│
└── database/migrations/
    ├── 001_example.sql
    └── ...
```

---

## Katmanlar

```
main.go
  ├── config     →  Yapılandırma okuma
  ├── utils      →  Yardımcı araçlar (bağımsız)
  ├── cache      →  Önbellek arayüzü ve servisi (Bu projede Redis tercih edildi)
  ├── database   →  Veritabanı arayüzü ve servisi (Bu projede MySQL tercih edildi)
  ├── storage    →  DB + Cache bir arada
  ├── events     →  Event sistemi + komut registry
  └── commands   →  Komutların kendisi
```

**Kural:** Bir katman sadece altındaki katmanı kullanır. `commands` → `utils`, `events`, `storage`. Hiçbir zaman tersi olmaz.
