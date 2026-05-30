package database

import (
	"VoidBot/config"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MysqlDB - MySQL veritabanı implementasyonu
// Database interface'ini implemente eder
type MysqlDB struct {
	db *sqlx.DB // sqlx bağlantısı (içinde connection pool var)
}

var DB *MysqlDB

// NewMySQL - yeni MySQL bağlantısı oluşturur
// Kullanım: db, err := NewMySQL("localhost", "3306", "root", "pass", "botdb")
func NewMySQL(cfg config.MySQLConfig) error {
	// Bağlantı string'i: parseTime=true ile time.Time otomatik dönüşür
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("mysql bağlantı hatası: %w", err)
	}

	// Connection pool ayarları
	db.SetMaxOpenConns(25)                  // Aynı anda max 25 bağlantı
	db.SetMaxIdleConns(10)                  // Boşta max 10 bağlantı beklet
	db.SetConnMaxIdleTime(10 * time.Minute) // 10 dk boş kalırsa kapat
	db.SetConnMaxLifetime(30 * time.Minute) // 30 dk sonra yenile

	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("mysql ping hatası: %w", err)
	}

	DB = &MysqlDB{db: db}
	return nil
}

// Query - çoklu satır dönen sorgular (SELECT * FROM ...)
func (m *MysqlDB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, query, args...)
}

// QueryRow - tek satır dönen sorgular (SELECT * FROM ... WHERE id = ?)
func (m *MysqlDB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return m.db.QueryRowContext(ctx, query, args...)
}
 
// Exec - INSERT, UPDATE, DELETE gibi değişiklik yapan sorgular
func (m *MysqlDB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}

// Select - çoklu satırları direkt struct slice'ına map eder
func (m *MysqlDB) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.db.SelectContext(ctx, dest, query, args...)
}

// Get - tek satırı direkt struct'a map eder
func (m *MysqlDB) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.db.GetContext(ctx, dest, query, args...)
}

// BeginTx - transaction başlatır (ya hep ya hiç)
func (m *MysqlDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Transaction, error) {
	tx, err := m.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &MysqlTx{tx: tx}, nil
}

// Close - bağlantı havuzunu kapatır
func (m *MysqlDB) Close() error {
	return m.db.Close()
}

// Ping - bağlantıyı test eder
func (m *MysqlDB) Ping() error {
	return m.db.Ping()
}

// MysqlTx - MySQL transaction implementasyonu
// Transaction interface'ini implemente eder
type MysqlTx struct {
	tx *sqlx.Tx // sqlx transaction'ı (Get/Select destekler)
}

// ExecContext - transaction içinde INSERT/UPDATE/DELETE
func (t *MysqlTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

// QueryContext - transaction içinde çoklu SELECT
func (t *MysqlTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

// QueryRowContext - transaction içinde tekli SELECT
func (t *MysqlTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

// SelectContext - transaction içinde çoklu satırı struct'a map eder
func (t *MysqlTx) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return t.tx.SelectContext(ctx, dest, query, args...)
}

// GetContext - transaction içinde tek satırı struct'a map eder
func (t *MysqlTx) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return t.tx.GetContext(ctx, dest, query, args...)
}

// Commit - transaction'ı onaylar (değişiklikleri kalıcı yapar)
func (t *MysqlTx) Commit() error {
	return t.tx.Commit()
}

// Rollback - transaction'ı geri alır (değişiklikleri iptal eder)
func (t *MysqlTx) Rollback() error {
	return t.tx.Rollback()
}
