package commands

import "github.com/bwmarrin/discordgo"

// Command - Her komutun uygulaması için bir arayüz. Her komut bu arayüzü uygulamalıdır.
type Command interface {
	Name() string
	Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
}

// Registry - Komutları saklamak için bir harita. Komut adı anahtar olarak kullanılır.
var registry = make(map[string]Command)

// Register - Yeni bir komutu kaydetmek için kullanılır. Komut adı benzersiz olmalıdır.
func Register(cmd Command) {
	registry[cmd.Name()] = cmd
}

// Handle - Gelen mesajları işler ve uygun komutu çalıştırır. Komut adı mesajın ilk kelimesi olarak varsayılır.
// args[0] komut adı, map'te varsa çalıştır
func Handle(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

	cmd, exist := registry[args[0]]
	if !exist {
		return // Komut yoksa sessiz kal
	}

	cmd.Execute(s, m, args)
}
