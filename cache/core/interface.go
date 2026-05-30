package cache

import "time"

// Cache - tüm cache işlemlerinin ortak arayüzü
type Cache interface {
	// Get - anahtara ait değeri getir, yoksa "", nil döner
	Get(key string) (string, error)

	// Set - değer kaydeder, expTime=0 ise süresiz
	Set(key string, value interface{}, expTime time.Duration) (bool, error)

	// Delete - anahtarı siler
	Delete(key string) error

	// Exists - anahtar var mı kontrol eder
	Exists(key string) (bool, error)

	// SetIfNotExists - anahtar yoksa kaydeder (cooldown için birebir)
	SetIfNotExists(key string, value interface{}, expTime time.Duration) (bool, error)

	// Incr/Decr - sayısal değeri 1 arttırır/azaltır
	Incr(key string) (int64, error)
	Decr(key string) (int64, error)

	// Expire - var olan anahtarın süresini değiştirir
	ExpireSet(key string, expTime time.Duration) error

	// Close/Ping - bağlantı yönetimi
	Close() error
	Ping() error
}
