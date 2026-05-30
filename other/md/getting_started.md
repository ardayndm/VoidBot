# Kurulum

## Gereksinimler

- Go 1.21+
- MySQL 8.0+
- Redis 7.0+
- Discord Bot Token

---

## 1. Repoyu klonla

```bash
git clone https://github.com/ardayndm/VoidBot.git
cd VoidBot
```

## 2. Bağımlılıkları yükle

```bash
go mod tidy
```

## 3. `.env` dosyasını oluştur

`.env.example` dosyasını kopyala ve doldur:

```bash
cp .env.example .env
```

```env
BOT_TOKEN=discord_bot_tokenin

MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=şifren
MYSQL_DATABASE=voidbot

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

LOG_LEVEL=ALL
```

## 4. Discord Developer Portal

1. [discord.com/developers](https://discord.com/developers/applications) adresine git
2. Uygulamanı seç → **Bot** sekmesi
3. **Message Content Intent** ✓ aç
4. **Server Members Intent** ✓ aç

## 5. Botu başlat

```bash
go run .
```

---

## Log Seviyeleri

| Değer     | Açıklama         |
| --------- | ---------------- |
| `ALL`     | Her şeyi yaz     |
| `INFO`    | Bilgi mesajları  |
| `WARNING` | Uyarılar         |
| `ERROR`   | Sadece hatalar   |
| `OK`      | Sadece başarılar |
| `NONE`    | Hiçbir şey yazma |
