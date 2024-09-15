package sqlx

import (
	"context"
	"fmt"

	"book-api/internal/domain"
	"book-api/internal/repository"
	"book-api/internal/transaction"

	"github.com/jmoiron/sqlx"
)

type SQLXBookRepository struct {
	DB *sqlx.DB
}

func NewSQLXBookRepository(db *sqlx.DB) repository.BookRepository {
	return &SQLXBookRepository{DB: db}
}

func NewBookPostgres(db *sqlx.DB) repository.BookRepository {
	return &SQLXBookRepository{DB: db}
}

const BookTable = "book"

func (r *SQLXBookRepository) Create(ctx context.Context, sqlxTx transaction.Transaction, book *domain.BookEntity) error {
	query := fmt.Sprintf(`INSERT INTO %s (name, description, cover, page_count, rate, author_id) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, BookTable)

	row := sqlxTx.QueryRowContext(ctx, query, book.BookInfo.Name, book.BookInfo.Description, book.BookInfo.Cover, book.BookInfo.PageCount, book.BookInfo.Rate, book.BookInfo.Author.Id)
	if err := row.Scan(&book.Id); err != nil {
		return fmt.Errorf("failed to insert book: %w", err)
	}

	return nil
}

func (r *SQLXBookRepository) Update(ctx context.Context, sqlxTx transaction.Transaction, book *domain.BookEntity) error {
	query := fmt.Sprintf(`UPDATE %s SET description = $2, cover = $3, page_count = $4, rate = $5, author_id = $6
							WHERE name = $1
							RETURNING id`, BookTable)

	row := sqlxTx.QueryRowContext(ctx, query, book.BookInfo.Name, book.BookInfo.Description, book.BookInfo.Cover, book.BookInfo.PageCount, book.BookInfo.Rate, book.BookInfo.Author.Id)
	if err := row.Scan(&book.Id); err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	return nil
}

func (r *SQLXBookRepository) GetByIds(ctx context.Context, ids []int) ([]domain.BookEntity, error) {
	var books []domain.BookEntity
	query := fmt.Sprintf(`
		SELECT 
			b.id as book_id, 
			b.cover, 
			b.description, 
			b.name as book_name, 
			b.page_count, 
			b.rate, 
			a.id as author_id, 
			a.name as author_name
		FROM %s AS b
		LEFT JOIN author AS a ON b.author_id = a.id
		WHERE b.id IN (?)`, BookTable)

	// TODO: replace above author table
	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to build book ID query: %w", err)
	}

	query = r.DB.Rebind(query)
	rows, err := r.DB.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select books: %w", err)
	}
	defer rows.Close()

	bookMap := make(map[int]*domain.BookEntity)

	for rows.Next() {
		var (
			bookId     int
			book       domain.BookEntity
			bookInfo   domain.BookInfo
			bookAuthor domain.BookAuthor
		)
		if err := rows.Scan(&bookId, &bookInfo.Cover, &bookInfo.Description, &bookInfo.Name, &bookInfo.PageCount, &bookInfo.Rate, &bookAuthor.Id, &bookAuthor.Name); err != nil {
			return nil, fmt.Errorf("failed to scan book data: %w", err)
		}

		book.Id = bookId
		book.BookInfo = bookInfo
		book.BookInfo.Author = &bookAuthor

		if existingBook, exists := bookMap[bookId]; exists {
			existingBook.BookInfo.Author = &bookAuthor
		} else {
			bookMap[bookId] = &book
		}
	}

	for _, book := range bookMap {
		books = append(books, *book)
	}

	return books, nil
}

func (r *SQLXBookRepository) GetAll(ctx context.Context) ([]domain.BookEntity, error) {
	var books []domain.BookEntity
	query := fmt.Sprintf(`
		SELECT 
			b.id as book_id, 
			b.cover, 
			b.description, 
			b.name as book_name, 
			b.page_count, 
			b.rate, 
			a.id as author_id, 
			a.name as author_name
		FROM %s AS b
		LEFT JOIN author AS a ON b.author_id = a.id;`, BookTable)

	rows, err := r.DB.QueryxContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select books: %w", err)
	}
	defer rows.Close()

	bookMap := make(map[int]*domain.BookEntity)

	for rows.Next() {
		var (
			bookId     int
			book       domain.BookEntity
			bookInfo   domain.BookInfo
			bookAuthor domain.BookAuthor
		)
		if err := rows.Scan(&bookId, &bookInfo.Cover, &bookInfo.Description, &bookInfo.Name, &bookInfo.PageCount, &bookInfo.Rate, &bookAuthor.Id, &bookAuthor.Name); err != nil {
			return nil, fmt.Errorf("failed to scan book data: %w", err)
		}

		book.Id = bookId
		book.BookInfo = bookInfo
		book.BookInfo.Author = &bookAuthor

		if existingBook, exists := bookMap[bookId]; exists {
			existingBook.BookInfo.Author = &bookAuthor
		} else {
			bookMap[bookId] = &book
		}
	}

	for _, book := range bookMap {
		books = append(books, *book)
	}

	return books, nil
}

func (r *SQLXBookRepository) Delete(ctx context.Context, sqlxTx transaction.Transaction, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", BookTable)

	_, err := sqlxTx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	return nil
}
