package repository

import (
	"context"
	"errors"
	"golibrary/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthorRepository struct {
	pool *pgxpool.Pool
}

func NewAuthorRepository(pool *pgxpool.Pool) *AuthorRepository {
	return &AuthorRepository{pool: pool}
}

// Create repositorys
func (a *AuthorRepository) CreateNewAuthors(ctx context.Context, author *models.Author) error {
	query := `
		INSERT INTO authors (username,email,password)
		VALUES ($1,$2,$3)
		RETURNING id,created_at
		`

	// Мы используем QueryRow, так как RETURNING возвращает только одну строчку
	err := a.pool.QueryRow(ctx, query, author.Username, author.Email, author.Password).Scan(&author.Id, &author.Created_at)
	if err != nil {
		return err
	}
	return nil
}

// Get repository
func (a *AuthorRepository) GetByIdAuthors(ctx context.Context, id uint64) (*models.Author, error) {
	query :=
		`
	SELECT id,username,email,created_at FROM authors WHERE id = $1
		`
	var author models.Author
	err := a.pool.QueryRow(ctx, query, id).Scan(
		&author.Id,
		&author.Username,
		&author.Email,
		&author.Created_at,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}
	return &author, nil
}

func (a AuthorRepository) GetAllAuthors(ctx context.Context) ([]models.Author, error) {
	query :=
		`
	SELECT * FROM authors
		`
	rows, err := a.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	authors, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Author])
	if err != nil {
		return nil, err
	}
	return authors, nil
}

func (a AuthorRepository) UpdateUsernameAuthors(ctx context.Context, id uint64, dto models.UpdateUsernameDTO) error {
	query :=
		`
	UPDATE authors SET username = $1 WHERE id = $2
	`

	cmdTag, err := a.pool.Exec(ctx, query, dto.Username, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return models.ErrUserNotFound
	}
	return nil
}

func (a *AuthorRepository) DeleteByIdAuthor(ctx context.Context, id uint64) error {
	query :=
		`
	DELETE FROM authors WHERE id = $1
		`
	cmdTag, err := a.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return models.ErrUserNotFound
	}
	return nil
}
