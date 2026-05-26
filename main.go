package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"voidbot/config"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Herşeyin başladığı yer.

	cfg := config.Load()

	dg, err := discordgo.New("Bot " + cfg.BotToken)

	if err != nil {
		log.Fatal("Discord oturumu oluşturulamadı: ", err)
	}

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
