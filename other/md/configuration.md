# Yapılandırma

---

## .env

Gizli ve ortama özel değerler burada tutulur. Asla Git'e gönderilmez.

```env
# Discord
BOT_TOKEN=your_discord_bot_token

# MySQL
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=
MYSQL_DATABASE=your_bot_database_name

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Log
LOG_LEVEL=ALL
```

---

## config/core/bot.yaml

Bot'un genel ayarları.

```yaml
name: "VoidBot"
version: "1.0.0"
author_id: "discord_kullanici_id"
prefix: "v!"
lang: "tr"
```

| Alan        | Açıklama                               |
| ----------- | -------------------------------------- |
| `name`      | Bot adı (embed footer'larında görünür) |
| `version`   | Bot versiyonu                          |
| `author_id` | Geliştirici Discord ID'si              |
| `prefix`    | Prefix komut eki (`v!warn` için `v!`)  |
| `lang`      | Aktif dil (`tr`, `en`)                 |

---

## config/core/color.yaml

Embed renkleri. Hex formatında girilir.

```yaml
success: "#00ff4c"
warning: "#ffbf00"
error: "#ff0011"
info: "#8f02b3"
```

---

## Veritabanı Migration

`database/migrations/` klasörüne `.sql` dosyaları ekle. Bot açılışında otomatik çalışır, bir kere çalışan tekrar çalışmaz.

```
database/migrations/
├── 001_example.sql
└── ...
```

Dosya adı sıralı olmalı — `001_`, `002_` gibi. Bot alfabetik sıraya göre çalıştırır.
