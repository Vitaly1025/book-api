package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"book-api/internal/domain"
	"book-api/internal/mapper"
	"book-api/internal/utils"
	dto "book-api/pkg/dto/http"

	"github.com/go-playground/validator/v10"
)

//go:generate mockgen -source=update_book.go -destination=mocks/update_book.go
type UpdateBookService interface {
	UpdateBook(ctx context.Context, req domain.BookEntity) (int, error)
}

type UpdateBookHandler struct {
	bookService UpdateBookService
	logger      *slog.Logger
	validator   *validator.Validate
}

func NewUpdateBookHandler(s UpdateBookService, l *slog.Logger) *UpdateBookHandler {
	return &UpdateBookHandler{bookService: s, logger: l, validator: validator.New()}
}

// @Summary Update book
// @Tags Book Operations
// @Description Update book
// @Accept text/json
// @Param request body models.Book true "Book"
// @Produce  json
// @Router /book [put]
// @Success  200  {int} resp
func (h *UpdateBookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	const op = "UpdateBook"
	log := h.logger.With(
		slog.String("operation", op),
	)

	ctx := context.Background()
	rw := utils.NewResponseWriter(w)
	rr := utils.NewRequestReader()

	var req dto.UpdateBookRequest
	err := rr.ParsePostRequest(r, &req)
	if err != nil {
		log.Error("cannot parse request", slog.String("error", err.Error()))
		err = rw.Error(http.StatusBadRequest, "invalid data")
		if err != nil {
			log.Error("cannot create error response", slog.String("error", err.Error()))
		}
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		log.Error("cannot pass validation", slog.String("error", err.Error()))
		err = rw.Error(http.StatusBadRequest, "data validation failed")
		if err != nil {
			log.Error("cannot create error response", slog.String("error", err.Error()))
		}
		return
	}

	query := mapper.ToDomainUpdateBook(req)

	resp, err := h.bookService.UpdateBook(ctx, query)
	if err != nil {
		log.Error("cannot update a book", slog.String("error", err.Error()))
		err = rw.Error(http.StatusBadRequest, "smth went wrong :(")
		if err != nil {
			log.Error("cannot create error response", slog.String("error", err.Error()))
		}
		return
	}

	err = rw.JSON(http.StatusOK, mapper.ToDTOCreateBook(resp))
	if err != nil {
		log.Error("cannot create json response", slog.String("error", err.Error()))
	}
}
