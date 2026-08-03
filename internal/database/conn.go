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
	// На Render файла .env нет, поэтому ошибки просто логируем без паники
	if err := godotenv.Load(); err != nil {
		log.Println("Информация: .env файл не найден, считываем переменные из ОС")
	}

	dbUrl := os.Getenv("DATABASE_URL")

	// Если DATABASE_URL пустой, пробуем собрать из отдельных переменных Render
	if dbUrl == "" {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")

		if dbHost != "" && dbUser != "" {
			dbUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
				dbUser, dbPass, dbHost, dbPort, dbName)
		}
	}

	return dbUrl
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
