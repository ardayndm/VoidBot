package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Presence - YAML'deki stringe göre discordgo tipi taslağı
var ActivityTypeMap = map[string]discordgo.ActivityType{
	"playing":   discordgo.ActivityTypeGame,
	"streaming": discordgo.ActivityTypeStreaming,
	"listening": discordgo.ActivityTypeListening,
	"watching":  discordgo.ActivityTypeWatching,
	"competing": discordgo.ActivityTypeCompeting,
}

// Presence - YAML'den okunan yapı
type PresenceConfig struct {
	Activities []struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
		URL  string `yaml:"url"`
	} `yaml:"activities"`
	Status   string `yaml:"status"`
	Interval int    `yaml:"interval"` // saniye
}

// Presence - presence.yaml dosyasını yükler
func loadPresenceConfig(lang string) (*PresenceConfig, error) {
	path := fmt.Sprintf("locale/%s/presence.yaml", lang)

	var cfg PresenceConfig
	if err := ReadYaml(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Presence - Bot durumunu rastgele günceller
func StartRandomPresence(s *discordgo.Session, botName, botLang string) {
	// Dil ayarını al
	lang := botLang
	if lang == "" {
		lang = "tr"
	}

	// Config'i yükle
	cfg, err := loadPresenceConfig(lang)
	if err != nil {
		Logger(WARNING, fmt.Sprintf("Presence config yüklenemedi: %v, varsayılan kullanılıyor", err))
		cfg = getDefaultPresenceConfig()
	}

	if len(cfg.Activities) == 0 {
		Logger(WARNING, "Presence aktivitesi bulunamadı")
		return
	}

	// Placeholder'ları doldur
	for i := range cfg.Activities {
		cfg.Activities[i].Name = FormatKeys(cfg.Activities[i].Name, map[string]string{
			"bot_name": botName,
		})
	}

	// Durum
	status := cfg.Status
	if status == "" {
		status = "online"
	}

	// Aralık
	interval := time.Duration(cfg.Interval) * time.Second
	if interval <= 0 {
		interval = 10 * time.Second
	}

	ticker := time.NewTicker(interval)

	Logger(OK, "Aktivite değiştirici başlatıldı.")
	go func() {
		for range ticker.C {
			// Rastgele aktivite seç
			act := cfg.Activities[rand.Intn(len(cfg.Activities))]

			activityType, ok := ActivityTypeMap[act.Type]
			if !ok {
				activityType = discordgo.ActivityTypeGame
			}

			activity := &discordgo.Activity{
				Name: act.Name,
				Type: activityType,
				URL:  act.URL,
			}

			err := s.UpdateStatusComplex(discordgo.UpdateStatusData{
				Status:     status,
				Activities: []*discordgo.Activity{activity},
			})

			if err != nil {
				Logger(WARNING, fmt.Sprintf("Durum güncellenemedi: %v", err))
			}
		}
	}()
}

// Presence - Varsayılan config (YAML okunamazsa)
func getDefaultPresenceConfig() *PresenceConfig {
	return &PresenceConfig{
		Activities: []struct {
			Name string `yaml:"name"`
			Type string `yaml:"type"`
			URL  string `yaml:"url"`
		}{
			{Name: "v!yardım", Type: "listening", URL: ""},
			{Name: "{bot_name} v-1.0.0", Type: "playing", URL: ""},
			{Name: "{bot_name} | v!yardım", Type: "watching", URL: ""},
		},
		Status:   "online",
		Interval: 10,
	}
}
