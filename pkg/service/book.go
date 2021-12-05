package service

import (
	models "book-api/pkg/model"
	"book-api/pkg/repository"
)

type BookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) *BookService{
	return &BookService{repo: repo}
}

func (s *BookService) CreateBook(req models.BookRequest) (*models.SimpleResponse, error){
	var resp models.SimpleResponse
	id, err := s.repo.CreateBook(req);
	if err != nil{
		return &resp, err
	}
	
	resp = models.SimpleResponse{
		Id: id,
	}

	return &resp, err;
}

func (s *BookService) UpdateBook(req models.Book) (*models.SimpleResponse, error){
	var resp models.SimpleResponse

	id, err := s.repo.UpdateBook(req)
	if err != nil{
		return &resp, err
	}
	
	resp = models.SimpleResponse{
		Id: id,
	}

	return &resp, err;
}

func (s *BookService) GetBookById(id int) (*models.Book, error){
	resp, err := s.repo.GetBook(id);
	return resp, err;
}

func (s *BookService) GetBooks(bookname string, genre string) (*[]models.Book, error){
	var resp *[]models.Book 
	var err error
	if len(bookname) > 0 && len(genre) >0 {
		resp, err = s.repo.GetBooksWithGeneralFilter(bookname, genre)
	} else if len(bookname) > 0 {
		resp, err = s.repo.GetBooksWithNameFilter(bookname)
	} else if len(genre) > 0 {
		resp, err = s.repo.GetBooksWithGenreFilter(genre)
	} else {
		resp, err = s.repo.GetAllBook()
	}
	
	if len(*resp) == 0 {
		return &[]models.Book{}, err
	} else{
		return resp, err;
	}
}

func (s *BookService) DeleteBook(id int) error{
	s.repo.DeleteBook(id)
	return nil
}