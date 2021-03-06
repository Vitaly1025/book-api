package validators

import (
	models "book-api/model"

	"github.com/go-playground/validator/v10"
)

func ValidateRequest(p models.BookRequest) error {
	validate := validator.New()
	return validate.Struct(p)
}
