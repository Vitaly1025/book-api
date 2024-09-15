package mapper

import (
	"book-api/internal/domain"
	dto "book-api/pkg/dto/http"
)

func ToDomainAuthor(name string) *domain.BookAuthor {
	return &domain.BookAuthor{Name: name}
}

func ToDomainGenres(genreIDs []int) []domain.BookGenre {
	genres := make([]domain.BookGenre, len(genreIDs))
	for i, id := range genreIDs {
		genres[i] = domain.BookGenre{Id: id}
	}
	return genres
}

func ToDomainCreateBook(req dto.CreateBookRequest) domain.BookEntity {
	return domain.BookEntity{
		BookInfo: domain.BookInfo{
			Author:      ToDomainAuthor(req.AuthorName),
			Cover:       req.Cover,
			Description: req.Description,
			Name:        req.Name,
			PageCount:   req.PageCount,
			Rate:        req.Rate,
		},
		Genres: ToDomainGenres(req.Genres),
	}
}

func ToDomainUpdateBook(req dto.UpdateBookRequest) domain.BookEntity {
	return domain.BookEntity{
		BookInfo: domain.BookInfo{
			Name:        req.Name,
			Author:      ToDomainAuthor(req.AuthorName),
			Cover:       req.Cover,
			Description: req.Description,
			PageCount:   req.PageCount,
			Rate:        req.Rate,
		},
		Genres: ToDomainGenres(req.Genres),
	}
}

func ToDTOCreateBook(id int) dto.SimpleResponse {
	return dto.SimpleResponse{
		ID: id,
	}
}
