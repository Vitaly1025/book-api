package service

import (
	"context"
	"fmt"

	"book-api/internal/domain"
	"book-api/internal/repository"
	transaction "book-api/internal/transaction"
)

type bookService struct {
	bookRepo   repository.BookRepository
	authorRepo repository.AuthorRepository
	genreRepo  repository.GenreRepository
	txManager  transaction.TransactionManager
}

func NewBookService(
	bookRepo repository.BookRepository,
	authorRepo repository.AuthorRepository,
	genreRepo repository.GenreRepository,
	txManager transaction.TransactionManager,
) BookService {
	return &bookService{bookRepo: bookRepo, authorRepo: authorRepo, genreRepo: genreRepo, txManager: txManager}
}

func (s *bookService) CreateBook(ctx context.Context, book domain.BookEntity) (int, error) {
	err := s.txManager.RunInTransaction(ctx, func(tx transaction.Transaction) error {
		existingAuthor, err := s.authorRepo.GetByName(ctx, book.BookInfo.Author.Name)
		if err != nil {
			return fmt.Errorf("failed to get author: %w", err)
		}

		if existingAuthor == nil {
			err = s.authorRepo.CreateOrUpdate(ctx, tx, book.BookInfo.Author)
			if err != nil {
				return fmt.Errorf("failed to create author: %w", err)
			}
		} else {
			book.BookInfo.Author = existingAuthor
		}

		err = s.bookRepo.Create(ctx, tx, &book)
		if err != nil {
			return fmt.Errorf("failed to create book: %w", err)
		}

		for _, genre := range book.Genres {
			err = s.genreRepo.UpsertGenreToBook(ctx, tx, book.Id, genre.Id)
			if err != nil {
				return fmt.Errorf("failed to add genre to book: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return -1, err
	}

	select {
	case <-ctx.Done():
		return -1, ctx.Err()
	default:
	}

	return book.Id, nil
}

func (s *bookService) UpdateBook(ctx context.Context, book domain.BookEntity) (int, error) {
	err := s.txManager.RunInTransaction(ctx, func(tx transaction.Transaction) error {
		existingAuthor, err := s.authorRepo.GetByName(ctx, book.BookInfo.Author.Name)
		if err != nil {
			return fmt.Errorf("failed to get author: %w", err)
		}

		if existingAuthor == nil {
			err = s.authorRepo.CreateOrUpdate(ctx, tx, book.BookInfo.Author)
			if err != nil {
				return fmt.Errorf("failed to update author: %w", err)
			}
		} else {
			book.BookInfo.Author = existingAuthor
		}

		err = s.bookRepo.Update(ctx, tx, &book)
		if err != nil {
			return fmt.Errorf("failed to upate book: %w", err)
		}

		for _, genre := range book.Genres {
			err = s.genreRepo.UpsertGenreToBook(ctx, tx, book.Id, genre.Id)
			if err != nil {
				return fmt.Errorf("failed to add genre to book: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return -1, err
	}

	select {
	case <-ctx.Done():
		return -1, ctx.Err()
	default:
	}

	return book.Id, nil
}

func (s *bookService) GetBookByIds(ctx context.Context, ids []int) ([]domain.BookEntity, error) {
	books, err := s.bookRepo.GetByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	// restore genres
	for i := range books {
		book := &books[i]
		bookGenres, err := s.genreRepo.GetByBookId(ctx, book.Id)
		if err != nil {
			return nil, err
		}
		book.Genres = bookGenres
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return books, nil
}

func (s *bookService) GetBooks(ctx context.Context) ([]domain.BookEntity, error) {
	books, err := s.bookRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// restore genres
	for i := range books {
		book := &books[i]
		bookGenres, err := s.genreRepo.GetByBookId(ctx, book.Id)
		if err != nil {
			return nil, err
		}
		book.Genres = bookGenres
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return books, nil
}

func (s *bookService) DeleteBook(ctx context.Context, id int) error {
	err := s.txManager.RunInTransaction(ctx, func(tx transaction.Transaction) error {
		err := s.genreRepo.DeleteBookGenre(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to remove bookgenre: %w", err)
		}

		err = s.bookRepo.Delete(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to remove book: %w", err)
		}

		return nil
	})

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return err
}
