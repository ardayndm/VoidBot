package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"voidbot/commands"
	"voidbot/commands/moderation"
	"voidbot/config"
	"voidbot/database"
	"voidbot/handlers"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Herşeyin başladığı yer.

	cfg := config.Load()

	initDb(cfg)
	registerCommands()

	dg, err := discordgo.New("Bot " + cfg.BotToken)

	if err != nil {
		log.Fatal("Discord oturumu oluşturulamadı: ", err)
	}

	dg.AddHandler(handlers.MessageReceived)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()

	if err != nil {
		log.Fatal("Discord oturumu açılamadı: ", err)
	}

	log.Println("VoidBot çalışıyor...")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	dg.Close()
}

// Veritabanı migrasyonlarını çalıştırır
func initDb(cfg *config.Config) {

	if err := database.Connect(cfg); err != nil {
		log.Fatal("Veritabanına bağlanılamadı: ", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatal("Veritabanı migrasyonu başarısız: ", err)
	}
}

func registerCommands() {
	// Burada tüm komutları kaydediyoruz. Yeni bir komut eklediğimizde buraya eklememiz gerekiyor.
	commands.Register(moderation.WarnCommand{})
}
