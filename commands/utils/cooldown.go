package utils

import (
	"sync"
	"time"
)

type CooldownManager struct {
	mu      sync.Mutex           // Aynı anda iki istek geldiğinde veri yarışını önlemek için mutex
	entries map[string]time.Time // "userID:command" -> son kullanım zamanı
}

// NewCooldownManager, Program başlayınca bir kere oluştur
func NewCooldownManager() *CooldownManager {
	return &CooldownManager{
		entries: make(map[string]time.Time),
	}
}

// Global bir Cooldown yöneticisi
var Cooldown = NewCooldownManager()

// Check - cooldown doldu mu ?
// true: kullanılabilir durumda , false: hala cooldownda - beklemeli.
func (c *CooldownManager) Check(userID, command string, duration time.Duration) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := userID + ":" + command
	last, exists := c.entries[key]

	if !exists {
		return true
	}
	return time.Since(last) >= duration
}

// Set - Kullanım zamanını kaydet
func (c *CooldownManager) Set(userID, command string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := userID + ":" + command
	c.entries[key] = time.Now()
}

// Remaining - Kalan süreyi döndürür, eğer cooldown dolmuşsa 0 döner
func (c *CooldownManager) Remaining(userID, command string, duration time.Duration) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := userID + ":" + command
	last, exists := c.entries[key]
	if !exists {
		return 0
	}
	remaining := duration - time.Since(last)

	if remaining < 0 {
		delete(c.entries, key) // Cooldown dolmuş, kaydı temizle
		return 0
	}

	return remaining
}
