package database

import (
	"context"
	"database/sql"
)

// Database - tüm veritabanı işlemlerinin ortak arayüzü
type Database interface {
	// Query - çoklu satır dönen SELECT için
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// QueryRow - tek satır dönen SELECT için
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row

	// Exec - INSERT, UPDATE, DELETE için
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// Select - çoklu satırı direkt struct slice'ına map eder
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Get - tek satırı direkt struct'a map eder
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// BeginTx - transaction başlatır (ya hep ya hiç)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Transaction, error)

	// Close - bağlantıyı kapatır
	Close() error

	// Ping - bağlantıyı test eder
	Ping() error
}

// Transaction - transaction içinde kullanılacak metodlar
type Transaction interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Commit() error   // Kaydet
	Rollback() error // Geri al
}
