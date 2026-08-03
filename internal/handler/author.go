package handler

import (
	"encoding/json"
	"errors"
	"golibrary/internal/models"
	"golibrary/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type AuthorHandler struct {
	service *service.AuthorService
}

func NewAuthorHandler(svc *service.AuthorService) *AuthorHandler {
	return &AuthorHandler{service: svc}
}

// POST-запросы
func (a *AuthorHandler) CreateAuthorHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.CreateAuthorDTO

	// Записываем данные в DTO, либо возвращаем ошибку
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// вызываем сервис и записываем все в структуру author
	author, err := a.service.CreateAuthor(r.Context(), dto)
	if err != nil {
		if errors.Is(err, models.ErrInvalidInput) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Выводим данные на экран
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(author)

}

// GET-запрсы
func (a *AuthorHandler) GetByIdAuthorHandler(w http.ResponseWriter, r *http.Request) {
	idSTR := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idSTR, 10, 64)
	if err != nil {
		http.Error(w, "invalid id param", http.StatusBadRequest)
		return
	}

	author, err := a.service.GetByIdService(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, models.ErrInvalidInput) {
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		log.Printf("ERROR in GetByIdAuthorHandler: %v\n", err)
		return
	}
	log.Printf("ERROR in GetByIdAuthorHandler: %v\n", err)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(author)
}

func (a *AuthorHandler) GetAllAuthorsHandler(w http.ResponseWriter, r *http.Request) {
	authors, err := a.service.GetAllAuthorsService(r.Context())

	if err != nil {
		log.Printf("[GetAllAuthorsHandler ERROR]: %v\n", err)
		http.Error(w, "internal server", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if authors == nil {
		authors = []models.Author{}
	}
	json.NewEncoder(w).Encode(authors)
}

// UPDATE-handlers
func (a *AuthorHandler) UpdateUsernameAuthorHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.UpdateUsernameDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if dto.Id <= 0 {
		http.Error(w, "invalid of missing authors", http.StatusBadRequest)
		return
	}

	err := a.service.UpdateUsernameAuthorsService(r.Context(), dto.Id, dto)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			http.Error(w, "Author not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "author has been updated username",
	})
}

func (a *AuthorHandler) DeleteByIdAuthorsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		http.Error(w, "Invalid id parametr", http.StatusBadRequest)
		return
	}

	if err := a.service.DeleteByIdAuthorsService(r.Context(), id); err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			http.Error(w, "author not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "succesfully",
	})
}
