package service

import (
	"context"
	"golibrary/internal/models"
	"golibrary/internal/repository"
	"strings"
)

type AuthorService struct {
	repo *repository.AuthorRepository
}

func NewAuthorService(repo *repository.AuthorRepository) *AuthorService {
	return &AuthorService{repo: repo}
}

func (a *AuthorService) CreateAuthor(ctx context.Context, dto models.CreateAuthorDTO) (*models.Author, error) {
	// Проверяем на пустоту строки
	if strings.TrimSpace(dto.Username) == "" || strings.TrimSpace(dto.Email) == "" || strings.TrimSpace(dto.Password) == "" {
		return nil, models.ErrInvalidInput
	}

	// Создаем экземпляр
	author := models.Author{
		Username: strings.TrimSpace(dto.Username),
		Email:    strings.TrimSpace(dto.Email),
		Password: strings.TrimSpace(dto.Password),
	}

	// Создаем SQL-запрос через репозиторий
	if err := a.repo.CreateNewAuthors(ctx, &author); err != nil {
		return nil, err
	}

	return &author, nil
}

// GETById authors
func (a *AuthorService) GetByIdService(ctx context.Context, id uint64) (*models.Author, error) {
	if id == 0 {
		return nil, models.ErrInvalidInput
	}
	return a.repo.GetByIdAuthors(ctx, id)
}

func (a *AuthorService) GetAllAuthorsService(ctx context.Context) ([]models.Author, error) {
	return a.repo.GetAllAuthors(ctx)
}

// UPDATE-запросы
func (a *AuthorService) UpdateUsernameAuthorsService(ctx context.Context, id uint64, dto models.UpdateUsernameDTO) error {
	if id <= 0 || strings.TrimSpace(dto.Username) == "" {
		return models.ErrInvalidInput
	}
	return a.repo.UpdateUsernameAuthors(ctx, id, dto)
}

// DELETE-запросы
func (a *AuthorService) DeleteByIdAuthorsService(ctx context.Context, id uint64) error {
	if id == 0 {
		return models.ErrInvalidInput
	}
	return a.repo.DeleteByIdAuthor(ctx, id)
}
