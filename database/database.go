package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"voidbot/config"
)

var DB *sql.DB

func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return err
	}

	// Bağlantıyı test et
	if err = db.Ping(); err != nil {
		return err
	}

	DB = db
	
	return nil
}
