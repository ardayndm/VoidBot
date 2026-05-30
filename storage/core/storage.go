package storage

import (
	cache "VoidBot/cache/core"
	database "VoidBot/database/core"
	"fmt"
)

type storage struct {
	DB    database.Database
	Cache cache.Cache
}

// Storage - Veritabanı ve cache işlemlerini yönetir
var storageManager *storage

// Storage - Esnek veritabanını verir
func GetDB() database.Database {
	if storageManager == nil || storageManager.DB == nil {
		return nil
	}
	return storageManager.DB
}

// Storage - Esnek cache'i verir
func GetCache() cache.Cache {
	if storageManager == nil || storageManager.Cache == nil {
		return nil
	}
	return storageManager.Cache
}

// Storage - Yeni storage oluşturur
func NewStorage(db database.Database, c cache.Cache, reqCache bool) error {
	if db == nil {
		return fmt.Errorf("veritabanı bağlantısı nil")
	}

	storageManager = &storage{
		DB:    db,
		Cache: c,
	}

	if _, err := IsDatabaseAvailable(); err != nil {
		return err
	}

	// Cache opsiyonel
	if reqCache {
		if _, err := IsCacheAvailable(); err != nil {
			return err
		}
	}

	return nil
}

// IsDatabaseAvailable - Veritabanı bağlantısını kontrol eder
// Bağlantı varsa true, yoksa false döner
func IsDatabaseAvailable() (available bool, e error) {
	if storageManager == nil || storageManager.DB == nil {
		return false, fmt.Errorf("storage veya database bağlantısı nil")
	}

	// Gerçek bağlantıyı test et (panic olursa false)
	defer func() {
		if r := recover(); r != nil {
			available = false
			e = fmt.Errorf("Veritabanı bağlantı hatası: %v", r)
		}
	}()

	// Ping dene
	if err := storageManager.DB.Ping(); err != nil {
		return false, fmt.Errorf("Veritabanı ping hatası: %w", err)
	}

	return true, nil
}

// IsCacheAvailable - Cache bağlantısını kontrol eder
func IsCacheAvailable() (available bool, e error) {
	if storageManager == nil || storageManager.Cache == nil {
		return false, fmt.Errorf("storage veya cache bağlantısı nil")
	}

	defer func() {
		if r := recover(); r != nil {
			available = false
			e = fmt.Errorf("Cache bağlantı hatası: %v", r)
		}
	}()

	if err := storageManager.Cache.Ping(); err != nil {

		return false, fmt.Errorf("Cache ping hatası: %w", err)
	}

	return true, nil
}
