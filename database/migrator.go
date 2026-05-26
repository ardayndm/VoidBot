package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"voidbot/commands/utils"
)

// Bu tablo, hangi migration'ların zaten çalıştırıldığını takip etmek için kullanılır. Böylece aynı migration'ı birden fazla kez çalıştırmaz.
// Bu tablo yoksa önce oluşturulmalıdır.
func createMigrationsTable() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id         INT AUTO_INCREMENT PRIMARY KEY,
			filename   VARCHAR(255) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

// Daha önce çalıştırılmış mı ?
func isMigrated(filename string) bool {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM migrations WHERE filename = ?", filename).Scan(&count)
	return count > 0
}

// Çalıştırıldı olarak işaretle
func markMigrated(filename string) error {
	_, err := DB.Exec("INSERT INTO migrations (filename) VALUES (?)", filename)
	return err
}

func Migrate() error {
	if err := createMigrationsTable(); err != nil {
		return fmt.Errorf("Migrations tablosu oluşturulamadı: %w", err)
	}

	// Dosyaları oku
	files, err := filepath.Glob("database/migrations/*.sql")
	if err != nil {
		return err
	}

	// Dosyaları sıralayarak çalıştır
	sort.Strings(files)

	for _, file := range files {
		filename := filepath.Base(file)

		// Zaten çalıştırıldıysa geç
		if isMigrated(filename) {
			utils.LogSuccess(fmt.Sprintf("%s zaten uygulandı, geçiliyor.\n", filename))
			continue
		}

		// SQL dosyasını oku
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("%s okunamadı: %w", filename, err)
		}

		// Birden fazla sorgu olabilir , noktalı virgülle ayır
		queries := strings.Split(string(content), ";")
		for _, query := range queries {
			query = strings.TrimSpace(query)
			if query == "" {
				continue
			}

			if _, err := DB.Exec(query); err != nil {
				return fmt.Errorf("%s çalıştırılamadı: %w", filename, err)
			}
		}

		// Çalıştırıldı olarak işaretle
		if err := markMigrated(filename); err != nil {
			return err
		}

		utils.LogSuccess(fmt.Sprintf("%s başarıyla uygulandı.\n", filename))

	}

	return nil
}
