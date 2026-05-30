package events

import (
	"github.com/bwmarrin/discordgo"
)

// Handler tipleri
type (
	MessageHandler     func(s *discordgo.Session, m *discordgo.MessageCreate)
	InteractionHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)
)

// EventBus — tüm handler'ları tutar ve dağıtır
type EventBus struct {
	message     []MessageHandler
	interaction []InteractionHandler
}

// Global bus — her yerden erişilir
var Bus = &EventBus{}

// ── KAYIT FONKSİYONLARI ──────────────────────────────────────────────

// OnMessage — MessageCreate event'i dinler
// Küfür filtresi, XP sistemi, prefix komutlar buraya kayıt olur
func (b *EventBus) OnMessage(h MessageHandler) {
	b.message = append(b.message, h)
}

// OnInteraction — InteractionCreate event'i dinler
// Slash komutlar buraya kayıt olur
func (b *EventBus) OnInteraction(h InteractionHandler) {
	b.interaction = append(b.interaction, h)
}

// ── DISPATCH FONKSİYONLARI ───────────────────────────────────────────
// Bu fonksiyonlar main.go'da Discord'a AddHandler ile bağlanır
// Gelen event'i kayıtlı tüm handler'lara sırayla iletir

// DispatchMessage — gelen mesajı tüm message handler'lara iletir
func (b *EventBus) DispatchMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, h := range b.message {
		h(s, m)
	}
}

// DispatchInteraction — gelen interaction'ı tüm interaction handler'lara iletir
func (b *EventBus) DispatchInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, h := range b.interaction {
		h(s, i)
	}
}

// ── İLERİDE YENİ EVENT EKLEMEK İÇİN ─────────────────────────────────
// Örnek: Ses kanalı desteği eklemek istersen:
//
// 1. Handler tipini ekle:
//    VoiceHandler func(s *discordgo.Session, v *discordgo.VoiceStateUpdate)
//
// 2. EventBus'a slice ekle:
//    voice []VoiceHandler
//
// 3. Kayıt fonksiyonu ekle:
//    func (b *EventBus) OnVoice(h VoiceHandler) { b.voice = append(b.voice, h) }
//
// 4. Dispatch fonksiyonu ekle:
//    func (b *EventBus) DispatchVoice(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
//        for _, h := range b.voice { h(s, v) }
//    }
//
// 5. main.go'ya bir satır ekle:
//    BotSession.AddHandler(events.Bus.DispatchVoice)
