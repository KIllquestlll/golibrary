package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func Load() string {
	if err := godotenv.Load(); err != nil {
		log.Printf("Строка пустая или не найден .env", err)
	}
	DBUrl := os.Getenv("DATABASE_URL")
	if DBUrl == "" {
		// log.Fatal("Переменная пустая!")
	}
	return DBUrl
}

func NewPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("Ошибка парсинга:%s", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		return nil, fmt.Errorf("Ошибка создания пула:%s", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("База данных не отвечает!", err)
	}
	return pool, nil
}
