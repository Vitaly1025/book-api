package sqlx

import (
	"context"
	"fmt"

	"book-api/internal/domain"
	"book-api/internal/repository"
	"book-api/internal/transaction"

	"github.com/jmoiron/sqlx"
)

type SQLXAuthorRepository struct {
	DB *sqlx.DB
}

func NewSQLXAuthorRepository(db *sqlx.DB) repository.AuthorRepository {
	return &SQLXAuthorRepository{DB: db}
}

const AuthorTable = "author"

func (r *SQLXAuthorRepository) CreateOrUpdate(ctx context.Context, sqlxTx transaction.Transaction, author *domain.BookAuthor) error {
	query := fmt.Sprintf(`INSERT INTO %s (name) VALUES ($1) 
						  ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name;
						  RETURNING id`, AuthorTable)

	row := sqlxTx.QueryRowContext(ctx, query, author.Name)
	if err := row.Scan(&author.Id); err != nil {
		return fmt.Errorf("failed to upsert author: %w", err)
	}

	return nil
}

func (r *SQLXAuthorRepository) GetByName(ctx context.Context, name string) (*domain.BookAuthor, error) {
	query := fmt.Sprintf(`SELECT id, name FROM %s WHERE name = $1 LIMIT 1`, AuthorTable)

	var author domain.BookAuthor
	err := r.DB.Get(&author, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get author by name: %w", err)
	}

	return &author, nil
}
