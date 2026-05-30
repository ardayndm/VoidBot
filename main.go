package main

import (
	"VoidBot/cache"
	config "VoidBot/config/core"
	"VoidBot/database"
	events "VoidBot/events/core"
	storage "VoidBot/storage/core"
	utils "VoidBot/utils/core"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func initConfigs() {
	utils.Logger(utils.INFO, "Yapılandırmalar yükleniyor...")

	// Konfigürasyonu yükle
	if err := config.InitConfig(); err != nil {
		utils.Logger(utils.ERROR, fmt.Sprintf("Config yüklenemedi: %s", err.Error()))
		os.Exit(1)
	}

	// Renkleri başlat
	if err := utils.InitColors(); err != nil {
		utils.Logger(utils.WARNING, fmt.Sprintf("Renkler yüklenemedi , devam ediliyor...: %s", err.Error()))
	}

	// Dil ayarlarını yükle
	if err := utils.InitLocale(config.AppConfig.Bot.Lang); err != nil {
		utils.Logger(utils.ERROR, fmt.Sprintf("Dil dosyası yüklenemedi: %s", err.Error()))
		os.Exit(1)
	}

	utils.Logger(utils.OK, "Yapılandırmalar yüklendi.")
}

func initDB() {
	utils.Logger(utils.INFO, "Veritabanı yükleniyor...")

	// Veritabanını yükle
	if err := database.NewMySQL(config.GetMySQL()); err != nil {
		utils.Logger(utils.ERROR, fmt.Sprintf("Veritabanı hatası: %s", err.Error()))
		os.Exit(1)
	}

	if database.DB == nil {
		utils.Logger(utils.ERROR, "Veritabanı bağlantısı yok!")
		os.Exit(1)
	}

	utils.Logger(utils.OK, "Veritabanı yüklendi.")

}

func initCache() {
	utils.Logger(utils.INFO, "Cache yükleniyor...")

	// Cache'i yükle
	if err := cache.NewRedis(config.GetRedis()); err != nil {
		utils.Logger(utils.WARNING, fmt.Sprintf("Cache devre dışı...: %s", err.Error()))
	}

	if cache.Redis == nil {
		utils.Logger(utils.WARNING, "Cache bağlantısı yok, devam ediliyor...")
	} else {
		utils.Logger(utils.OK, "Cache yüklendi.")
	}
}

func initStorage() {
	utils.Logger(utils.INFO, "Storage kuruluyor...")
	// Esnek kullanım için interface üzerinden erişimi hazırla
	if err := storage.NewStorage(database.DB, cache.Redis, false); err != nil {
		utils.Logger(utils.ERROR, fmt.Sprintf("Storage kurulamadı: %s", err.Error()))
		os.Exit(1)
	}
	utils.Logger(utils.OK, "Storage kuruldu.")
}

func initDatabaseMigration() {
	utils.Logger(utils.INFO, "Veritabanı tabloları yükleniyor...")
	// Migration'ı çalıştır (DB Tablolarını hazırla)
	migration := storage.NewMigrationManager()
	if err := migration.Migrate(); err != nil {
		utils.Logger(utils.ERROR, fmt.Sprintf("Migration hatası: %v", err))
		os.Exit(1)
	}
	utils.Logger(utils.OK, "Veritabanı tabloları yüklendi.")
}

// Tüm başlangıç ayarlarını yapar
func init() {

	// .env dosyasını yükle
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf(".env dosyası yüklenemedi: %w", err))
	}

	// Uygulama Log Level bilgilerini yükle
	lvl := os.Getenv("LOG_LEVEL")
	utils.SetLogLevel(lvl)

	utils.Logger(utils.INFO, "Ön yükleme başlatılıyor...")

	initConfigs()
	// initDB()
	// initCache()
	// initStorage()
	// initDatabaseMigration()

	/*
		storage.GetDB().Ping()
		Örnek kullanma şekli

	*/

	/*
		Komut yükleme/kullanma işlemi
		commands, err := utils.LoadCommand("komut_adı")
	*/

	utils.Logger(utils.OK, "Önyükleme tamamlandı.")
}

var BotSession *discordgo.Session

func main() {
	utils.Logger(utils.INFO, "Bot başlatılıyor...")
	bot := config.GetBot()

	// Discord oturumu başlat
	var err error
	BotSession, err = discordgo.New("Bot " + bot.Token)
	if err != nil {
		utils.Logger(utils.ERROR, fmt.Sprintf("Discord bağlantı hatası: %v", err))
		os.Exit(1)
	}

	// ==========================  Event handler =======================
	// Handler'ları bus'a bağla
	events.Bus.OnMessage(events.HandleMessage)
	events.Bus.OnInteraction(events.HandleInteraction)

	// Discord'a bağla
	BotSession.AddHandler(events.Bus.DispatchMessage)
	BotSession.AddHandler(events.Bus.DispatchInteraction)

	// =================================================================

	// Botu aç
	if err := BotSession.Open(); err != nil {
		utils.Logger(utils.ERROR, fmt.Sprintf("Bot açılamadı: %v", err))
		os.Exit(1)
	}
	defer BotSession.Close()

	utils.Logger(utils.OK, fmt.Sprintf("%s çalışıyor! (Prefix: %s)", bot.Name, bot.Prefix))

	// Discord'a botun komutlarını bildir
	events.SyncSlashCommands(BotSession)

	// Graceful shutdown - Ctrl+C bekler
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	utils.Logger(utils.INFO, "Bot kapatılıyor...")
	shutdown()
}

// shutdown - Kaynakları temizler
func shutdown() {
	if database.DB != nil {
		database.DB.Close()
	}
	if cache.Redis != nil {
		cache.Redis.Close()
	}

	utils.Logger(utils.OK, "Bot başarıyla kapatıldı")
}
