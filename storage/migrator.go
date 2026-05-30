package storage

import (
	"VoidBot/utils"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// MigrationManager - Migration işlemlerini yönetir
type MigrationManager struct {
	store *storage // İsim değişti! (storage yerine store)
}

// NewMigrationManager - Yeni migration manager oluşturur
func NewMigrationManager() *MigrationManager {
	return &MigrationManager{
		store: storageManager,
	}
}

// createMigrationsTable - Migration takip tablosunu oluşturur
func (m *MigrationManager) createMigrationsTable() error {
	ctx, cancel := utils.Ctx(5)
	defer cancel()
	_, err := m.store.DB.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS migrations (
			id         INT AUTO_INCREMENT PRIMARY KEY,
			filename   VARCHAR(255) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

// isMigrated - Migration daha önce çalıştırıldı mı?
func (m *MigrationManager) isMigrated(filename string) bool {
	ctx, cancel := utils.Ctx(5)
	defer cancel()
	var count int
	err := m.store.DB.QueryRow(ctx, "SELECT COUNT(*) FROM migrations WHERE filename = ?", filename).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

// markMigrated - Migration'ı çalıştırıldı olarak işaretler
func (m *MigrationManager) markMigrated(filename string) error {
	ctx, cancel := utils.Ctx(5)
	defer cancel()
	_, err := m.store.DB.Exec(ctx, "INSERT INTO migrations (filename) VALUES (?)", filename)
	return err
}

// Migrate - Tüm migration'ları çalıştırır
func (m *MigrationManager) Migrate() error {
	utils.Logger(utils.INFO, "Migration başlatılıyor...")

	if ok, err := IsDatabaseAvailable(); !ok {
		return err
	}

	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("migration tablosu oluşturulamadı: %w", err)
	}

	files, err := filepath.Glob("database/migrations/*.sql")
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {
		filename := filepath.Base(file)

		if m.isMigrated(filename) {
			utils.Logger(utils.INFO, fmt.Sprintf("%s zaten uygulandı, atlanıyor", filename))
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("%s okunamadı: %w", filename, err)
		}

		utils.Logger(utils.INFO, fmt.Sprintf("%s uygulanıyor...", filename))

		queries := strings.Split(string(content), ";")
		var execErr error
		for _, query := range queries {
			query = strings.TrimSpace(query)
			if query == "" {
				continue
			}

			func() {
				ctx, cancel := utils.Ctx(5)
				defer cancel()
				_, execErr = m.store.DB.Exec(ctx, query)
			}()

			if execErr != nil {
				return fmt.Errorf("%s çalıştırılamadı: %w\nSorgu: %s", filename, execErr, query)
			}
		}

		if err := m.markMigrated(filename); err != nil {
			return err
		}

		utils.Logger(utils.OK, fmt.Sprintf("%s başarıyla uygulandı", filename))
	}

	utils.Logger(utils.OK, "Tüm migration'lar tamamlandı!")
	return nil
}
