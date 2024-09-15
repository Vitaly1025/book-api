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

//go:generate mockgen -source=create_book.go -destination=mocks/create_book.go
type CreateBookService interface {
	CreateBook(ctx context.Context, req domain.BookEntity) (int, error)
}

type CreateBookHandler struct {
	bs        CreateBookService
	logger    *slog.Logger
	validator *validator.Validate
}

func NewCreateBookHandler(s CreateBookService, l *slog.Logger) *CreateBookHandler {
	return &CreateBookHandler{bs: s, logger: l, validator: validator.New()}
}

const op = "CreateBook"

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book entry in the system
// @Tags Books
// @Accept json
// @Produce json
// @Param request body dto.CreateBookRequest true "Create Book Request"
// @Success 200 {object} dto.SimpleResponse "Book created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid data or validation failed"
// @Router /books [post]
func (h *CreateBookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(
		slog.String("operation", op),
	)
	ctx := context.Background()
	rw := utils.NewResponseWriter(w)
	rr := utils.NewRequestReader()

	var req dto.CreateBookRequest
	err := rr.ParsePostRequest(r, &req)
	if err != nil {
		log.Error("cannot parse post request", slog.String("error", err.Error()))
		err = rw.Error(http.StatusBadRequest, "invalid data")
		if err != nil {
			log.Error("cannot create error response", slog.String("error", err.Error()))
		}
		return
	}

	err = h.validator.Struct(&req)
	if err != nil {
		log.Error("cannot validate post request", slog.String("error", err.Error()))
		err = rw.Error(http.StatusBadRequest, "data validation failed")
		if err != nil {
			log.Error("cannot create error response", slog.String("error", err.Error()))
		}
		return
	}

	query := mapper.ToDomainCreateBook(req)

	resp, err := h.bs.CreateBook(ctx, query)
	if err != nil {
		log.Error("cannot create book", slog.String("error", err.Error()))
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
