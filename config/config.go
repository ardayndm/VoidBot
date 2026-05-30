package config

import (
	"VoidBot/utils"
	"fmt"
	"os"
	"strconv"
)

// BotConfig - Discord bot ayarları
type BotConfig struct {
	Token    string
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	AuthorID string `yaml:"author_id"`
	Prefix   string `yaml:"prefix"`
	Lang     string `yaml:"lang"`
}

// MySQLConfig - MySQL veritabanı ayarları
type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string // DB Name
}

// RedisConfig - Redis cache ayarları
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// Config - Tüm konfigürasyonları içeren ana yapı
type Config struct {
	Bot      BotConfig
	MySQL    MySQLConfig
	Redis    RedisConfig
	LogLevel utils.LogLevel
}

// Global ayarlar
var AppConfig *Config

// InitConfig - Tüm konfigürasyonları .env dosyasından yükler
func InitConfig() error {
	cfg := &Config{}

	// Bot ayarlarını yükle
	if err := loadBotConfig(cfg); err != nil {
		return err
	}

	// MySQL ayarlarını yükle
	if err := loadMySQLConfig(cfg); err != nil {
		return err
	}

	// Redis ayarlarını yükle
	if err := loadRedisConfig(cfg); err != nil {
		return err
	}

	AppConfig = cfg
	return nil
}

// loadBotConfig - Bot ayarlarını yükler
func loadBotConfig(cfg *Config) error {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return fmt.Errorf("BOT_TOKEN bulunamadı")
	}

	if err := utils.ReadYaml("config/bot.yaml", &cfg.Bot); err != nil {
		return err
	}

	cfg.Bot.Token = token

	if cfg.Bot.Lang == "" {
		cfg.Bot.Lang = "tr"
	}

	return nil
}

// loadMySQLConfig - MySQL ayarlarını yükler
func loadMySQLConfig(cfg *Config) error {
	cfg.MySQL.Host = os.Getenv("MYSQL_HOST")
	cfg.MySQL.Port = os.Getenv("MYSQL_PORT")
	cfg.MySQL.User = os.Getenv("MYSQL_USER")
	cfg.MySQL.Password = os.Getenv("MYSQL_PASSWORD")
	cfg.MySQL.Database = os.Getenv("MYSQL_DATABASE")

	// Varsayılan port
	if cfg.MySQL.Port == "" {
		cfg.MySQL.Port = "3306"
	}

	// Zorunlu alanları kontrol et
	if cfg.MySQL.Host == "" {
		return fmt.Errorf("MYSQL_HOST bulunamadı")
	}
	if cfg.MySQL.User == "" {
		return fmt.Errorf("MYSQL_USER bulunamadı")
	}
	if cfg.MySQL.Database == "" {
		return fmt.Errorf("MYSQL_DATABASE bulunamadı")
	}

	return nil
}

// loadRedisConfig - Redis ayarlarını yükler
func loadRedisConfig(cfg *Config) error {
	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	cfg.Redis.Port = os.Getenv("REDIS_PORT")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")

	// Redis DB numarasını parse et
	dbStr := os.Getenv("REDIS_DB")
	if dbStr == "" {
		cfg.Redis.DB = 0 // Varsayılan
	} else {
		db, err := strconv.Atoi(dbStr)
		if err != nil {
			return fmt.Errorf("REDIS_DB geçersiz: %w", err)
		}
		cfg.Redis.DB = db
	}

	// Varsayılan port
	if cfg.Redis.Port == "" {
		cfg.Redis.Port = "6379"
	}

	// Varsayılan host
	if cfg.Redis.Host == "" {
		cfg.Redis.Host = "localhost"
	}

	return nil
}

// GetBot - Bot config'ini döndür
func GetBot() BotConfig {
	return AppConfig.Bot
}

// GetMySQL - MySQL config'ini döndür
func GetMySQL() MySQLConfig {
	return AppConfig.MySQL
}

// GetRedis - Redis config'ini döndür
func GetRedis() RedisConfig {
	return AppConfig.Redis
}
