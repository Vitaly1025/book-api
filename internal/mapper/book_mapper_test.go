package mapper_test

import (
	"testing"

	"book-api/internal/domain"
	"book-api/internal/mapper"
	dto "book-api/pkg/dto/http"

	"github.com/magiconair/properties/assert"
)

func TestToDomainAuthor(t *testing.T) {
	name := "Test Author"
	expected := &domain.BookAuthor{Name: name}

	result := mapper.ToDomainAuthor(name)
	assert.Equal(t, expected, result)
}

func TestToDomainGenres(t *testing.T) {
	genreIDs := []int{1, 2, 3}
	expected := []domain.BookGenre{
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}

	result := mapper.ToDomainGenres(genreIDs)
	assert.Equal(t, expected, result)
}

func TestToDomainCreateBook(t *testing.T) {
	req := dto.CreateBookRequest{
		Name:        "Test Book",
		AuthorName:  "Test Author",
		Description: "This is a test book",
		Genres:      []int{1, 2},
		Cover:       []byte("test cover"),
		PageCount:   100,
		Rate:        5,
	}

	expected := domain.BookEntity{
		BookInfo: domain.BookInfo{
			Name:        "Test Book",
			Author:      &domain.BookAuthor{Name: "Test Author"},
			Description: "This is a test book",
			Cover:       []byte("test cover"),
			PageCount:   100,
			Rate:        5,
		},
		Genres: []domain.BookGenre{
			{Id: 1},
			{Id: 2},
		},
	}

	result := mapper.ToDomainCreateBook(req)
	assert.Equal(t, expected, result)
}
