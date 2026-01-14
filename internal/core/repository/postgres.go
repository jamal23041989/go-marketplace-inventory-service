package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jamal23041989/go-marketplace-inventory-service/internal/core/config"
	_ "github.com/lib/pq"
)

const (
	MaxOpenConns    = 10
	MaxIdleConns    = 5
	MaxLifetimeConn = 30
)

func InitDB(cfg config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	for i := 1; i <= 5; i++ {
		if err := db.Ping(); err == nil {
			break
		}

		if i != 5 {
			time.Sleep(2 * time.Second)
			continue
		}

		return nil, err
	}

	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(MaxLifetimeConn) * time.Second)

	return db, nil
}
