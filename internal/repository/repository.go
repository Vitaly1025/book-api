package repository

import (
	"context"

	"book-api/internal/domain"
	"book-api/internal/transaction"
)

//go:generate mockgen -source=repository.go -destination=./mocks/repository.go
type BookRepository interface {
	Create(ctx context.Context, tx transaction.Transaction, book *domain.BookEntity) error
	Update(ctx context.Context, tx transaction.Transaction, book *domain.BookEntity) error
	// GetByPredicate(ctx context.Context, predicate string) ([]domain.BookEntity, error)
	GetByIds(ctx context.Context, ids []int) ([]domain.BookEntity, error)
	Delete(ctx context.Context, tx transaction.Transaction, id int) error
	GetAll(ctx context.Context) ([]domain.BookEntity, error)
}

//go:generate mockgen -source=repository.go -destination=./mocks/repository.go
type AuthorRepository interface {
	CreateOrUpdate(ctx context.Context, tx transaction.Transaction, author *domain.BookAuthor) error
	GetByName(ctx context.Context, name string) (*domain.BookAuthor, error)
}

//go:generate mockgen -source=repository.go -destination=./mocks/repository.go
type GenreRepository interface {
	DeleteBookGenre(ctx context.Context, tx transaction.Transaction, bookId int) error
	UpsertGenreToBook(ctx context.Context, tx transaction.Transaction, bookId, genreId int) error
	GetByBookId(ctx context.Context, bookId int) ([]domain.BookGenre, error)
}
