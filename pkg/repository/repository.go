package repository

import (
	models "book-api/pkg/model"

	"github.com/jmoiron/sqlx"
)

type BookRepository interface {
	CreateBook(book models.BookRequest) (int, error)
	UpdateBook(b models.Book) (int,error)
	GetBook(id int) (*models.Book, error)
	GetBooksWithNameFilter(name string) (*[]models.Book, error)
	GetBooksWithGenreFilter(name string) (*[]models.Book, error)
	GetAllBook() (*[]models.Book, error)
	GetBooksWithGeneralFilter(bookname string, genre string) (*[]models.Book, error)
	DeleteBook(id int) error
}

type Repository struct {
	BookRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		BookRepository: NewBookPostgres(db),
	}
}