package service

import (
	"context"

	"book-api/internal/domain"
)

type BookService interface {
	CreateBook(ctx context.Context, req domain.BookEntity) (int, error)
	UpdateBook(ctx context.Context, req domain.BookEntity) (int, error)
	GetBookByIds(ctx context.Context, ids []int) ([]domain.BookEntity, error)
	GetBooks(ctx context.Context) ([]domain.BookEntity, error)
	DeleteBook(ctx context.Context, id int) error
}
