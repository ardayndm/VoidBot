package config

// BotConfig, botun temel yapılandırma bilgilerini tutar.
// Sunucudan sunucuya değişebilecek bilgiler her şey burada
type BotConfig struct {
	Prefix  string
	Name    string
	Version string
}

func Default() *BotConfig {

	return &BotConfig{
		Prefix:  "!",
		Name:    "VoidBot",
		Version: "0.0.1",
	}
}

// Bot, uygulama genelinde kullanılacak tek bir BotConfig örneğidir.
var Bot = Default()
