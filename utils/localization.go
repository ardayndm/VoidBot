package utils

import "fmt"

// Localization - common.yaml içerisindeki title + messages
type MessageSection struct {
	Title    string            `yaml:"title"`
	Messages map[string]string `yaml:"messages"`
}

// Localization - common.yaml ile birebir eşleşir
type CommonLocale struct {
	Error   MessageSection    `yaml:"error"`
	Warning MessageSection    `yaml:"warning"`
	Success MessageSection    `yaml:"success"`
	Info    MessageSection    `yaml:"info"`
	Icons   map[string]string `yaml:"icons"`
}

// Localization -  Her komutun yaml yapısı (command_example.yaml patterni)
type CommandLocale struct {
	Titles   map[string]string `yaml:"titles"`
	Messages map[string]string `yaml:"messages"`
	Options  map[string]any    `yaml:"options"`
}

// Localization -
type locale struct {
	lang     string
	common   *CommonLocale
	commands map[string]*CommandLocale
}

var loc *locale

// Localization - Dil yapısını kurar.
func InitLocale(lang string) error {
	if lang == "" {
		lang = "tr"
	}

	common := &CommonLocale{}
	if err := ReadYaml(fmt.Sprintf("locale/%s/common.yaml", lang), common); err != nil {
		return fmt.Errorf("common.yaml okunamadı: %w", err)
	}

	loc = &locale{
		lang:     lang,
		common:   common,
		commands: make(map[string]*CommandLocale),
	}

	return nil
}

// Localization - common.yaml - Genel mesajlara erişim
func Common() *CommonLocale {
	return loc.common
}

// Localization - Aktif dili döndürür
func GetLanguage() string {
	return loc.lang
}

// Localization — Options map'inden int değer okur, bulamazsa fallback döner
func OptionInt(cmd *CommandLocale, key string, fallback int) int {
	val, ok := cmd.Options[key]
	if !ok {
		return fallback
	}
	n, ok := val.(int)
	if !ok {
		return fallback
	}
	return n
}

// Localization - Options map'inden bool değer okur, bulamazsa fallback döner
func OptionBool(cmd *CommandLocale, key string, fallback bool) bool {
	val, ok := cmd.Options[key]
	if !ok {
		return fallback
	}
	b, ok := val.(bool)
	if !ok {
		return fallback
	}
	return b
}

// Localization - Options map'inden string değer okur, bulamazsa fallback döner
func OptionString(cmd *CommandLocale, key string, fallback string) string {
	val, ok := cmd.Options[key]
	if !ok {
		return fallback
	}
	s, ok := val.(string)
	if !ok {
		return fallback
	}
	return s
}

// Localization - Komut yaml'ını lazy okur, cache'den döner
func LoadCommand(name string) (*CommandLocale, error) {
	// Eğer map de var ise direkt mapden döndür
	if cmd, ok := loc.commands[name]; ok {
		return cmd, nil
	}

	// mapde yok ise dosyayı bulup yüklemeye çalış
	cmd := &CommandLocale{}
	path := fmt.Sprintf("locale/%s/commands/%s.yaml", loc.lang, name)

	// Dosya bulunamazsa hata döndür
	if err := ReadYaml(path, cmd); err != nil {
		return nil, fmt.Errorf("%s.yaml okunamadı: %w", name, err)
	}

	// Dosyayı mape kaydet ve döndür
	loc.commands[name] = cmd
	return cmd, nil
}
