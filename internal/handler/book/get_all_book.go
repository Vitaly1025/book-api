package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"book-api/internal/domain"
	"book-api/internal/utils"
)

//go:generate mockgen -source=get_all_book.go -destination=mocks/get_all_book.go
type GetAllBookService interface {
	GetBooks(ctx context.Context) ([]domain.BookEntity, error)
}

type GetAllBookHandler struct {
	bookService GetAllBookService
	logger      *slog.Logger
}

func NewGetAllBookHandler(s GetAllBookService, l *slog.Logger) *GetAllBookHandler {
	return &GetAllBookHandler{bookService: s, logger: l}
}

// @Summary Get books
// @Tags Book Operations
// @Description This method return books
// @Accept text/json
// @Param predicate query string false "BookName"
// @Param genre query string false "GenreName"
// @Produce  json
// @Router /book [get]
func (h *GetAllBookHandler) GetAllBook(w http.ResponseWriter, r *http.Request) {
	const op = "GetAllBook"
	log := h.logger.With(
		slog.String("operation", op),
	)
	ctx := context.Background()
	rw := utils.NewResponseWriter(w)

	resp, err := h.bookService.GetBooks(ctx)
	if err != nil {
		log.Error("cannot get a book", slog.String("error", err.Error()))
		err = rw.Error(http.StatusBadRequest, "can't got a book")
		if err != nil {
			log.Error("cannot create error response", slog.String("error", err.Error()))
		}
		return
	}

	err = rw.JSON(http.StatusOK, resp)
	if err != nil {
		log.Error("cannot create json response", slog.String("error", err.Error()))
	}
}
