package service

import (
	models "book-api/model"
	"book-api/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type BookWork interface {
	CreateBook(req models.BookRequest) (*models.SimpleResponse, error)
	UpdateBook(req models.Book) (*models.SimpleResponse, error)
	GetBookById(id int) (*models.Book, error)
	GetBooks(bookname string, genre string) (*[]models.Book, error)
	DeleteBook(id int) error
}

type Service struct {
	BookWork
}

func NewService(repos *repository.Repository) *Service { 
	return &Service{
		BookWork: NewBookService(repos.BookRepository),
	} 
}