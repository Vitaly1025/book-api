package sqlx

import (
	"context"
	"fmt"

	"book-api/internal/domain"
	"book-api/internal/repository"
	"book-api/internal/transaction"

	"github.com/jmoiron/sqlx"
)

type SQLXAGenreRepository struct {
	DB *sqlx.DB
}

func NewSQLXAGenreRepository(db *sqlx.DB) repository.GenreRepository {
	return &SQLXAGenreRepository{DB: db}
}

const (
	BookGenreTable = "book_genre"
	GenreTable     = "genre"
)

func (r *SQLXAGenreRepository) UpsertGenreToBook(ctx context.Context, sqlxTx transaction.Transaction, bookID int, genreID int) error {
	query := fmt.Sprintf(`INSERT INTO %s (book_id, genre_id) 
							VALUES ($1, $2) 
							ON CONFLICT (book_id) 
							DO UPDATE SET genre_id = EXCLUDED.genre_id`, BookGenreTable)
	_, err := sqlxTx.ExecContext(ctx, query, bookID, genreID)
	if err != nil {
		return fmt.Errorf("failed to upsert genre to book: %w", err)
	}
	return nil
}

func (r *SQLXAGenreRepository) GetByBookId(ctx context.Context, bookId int) ([]domain.BookGenre, error) {
	var bookGenre []domain.BookGenre
	query := fmt.Sprintf(`SELECT G.id AS Id, G.name AS Name
							FROM %s AS BG 
							INNER JOIN %s AS G ON BG.genre_id = G.id
							WHERE BG.book_id = $1;`, BookGenreTable, GenreTable)
	query = r.DB.Rebind(query)
	err := r.DB.Select(&bookGenre, query, bookId)
	if err != nil {
		return bookGenre, fmt.Errorf("failed to select genrebook by book id: %w", err)
	}
	return bookGenre, nil
}

func (r *SQLXAGenreRepository) DeleteBookGenre(ctx context.Context, sqlxTx transaction.Transaction, bookId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE book_id = $1", BookGenreTable)

	_, err := sqlxTx.ExecContext(ctx, query, bookId)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	return nil
}
