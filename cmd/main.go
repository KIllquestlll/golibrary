package main

import (
	"context"
	"golibrary/internal/config"
	"golibrary/internal/database"
	"golibrary/internal/handler"
	"golibrary/internal/repository"
	"golibrary/internal/service"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Инцилизация конфига
	cfg := config.Load()
	ctx := context.Background()
	dbURL := database.Load()

	// Создания роутера
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	pool, err := database.NewPool(ctx, dbURL)
	if err != nil {
		log.Fatal("Ошика подлючения к бд!", err)
	}

	if err := database.InitTables(ctx, pool); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}
	defer pool.Close()

	// Helper handlers
	helper := handler.HelperHandler(pool)
	r.Get("/ping", helper.Ping)

	// author handler
	authorRepo := repository.NewAuthorRepository(pool)
	authorSvc := service.NewAuthorService(authorRepo)
	author := handler.NewAuthorHandler(authorSvc)

	r.Route("/author", func(r chi.Router) {
		r.Post("/create", author.CreateAuthorHandler)
		r.Get("/", author.GetAllAuthorsHandler)
		r.Get("/{id}", author.GetByIdAuthorHandler)
		r.Put("/update", author.UpdateUsernameAuthorHandler)
		r.Delete("/delete/{id}", author.DeleteByIdAuthorsHandler)
	})

	// Запуск сервера
	srv := config.New(cfg, r)
	log.Printf("Server has been started:%s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
