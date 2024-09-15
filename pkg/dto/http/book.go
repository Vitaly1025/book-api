package dto

type CreateBookRequest struct {
	Name        string  `json:"name" validate:"required"`
	AuthorName  string  `json:"author_name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Genres      []int   `json:"genres" validate:"required,dive,required"`
	Cover       []byte  `json:"cover" validate:"required"`
	PageCount   int     `json:"page_count" validate:"required,min=1"`
	Rate        float32 `json:"rate" validate:"min=0,max=10"`
}

type UpdateBookRequest struct {
	Name        string  `json:"name" validate:"required"`
	AuthorName  string  `json:"author_name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Genres      []int   `json:"genres" validate:"required,dive,required"`
	Cover       []byte  `json:"cover" validate:"required"`
	PageCount   int     `json:"page_count" validate:"required,min=1"`
	Rate        float32 `json:"rate" validate:"min=0,max=10"`
}

type BookResponse struct {
	ID          int      `json:"id"`
	AuthorName  string   `json:"author_name"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Genres      []string `json:"genres"`
	Cover       []byte   `json:"cover"`
	PageCount   int      `json:"page_count"`
	Rate        int      `json:"rate"`
}

type SimpleResponse struct {
	ID int `json:"id"`
}
