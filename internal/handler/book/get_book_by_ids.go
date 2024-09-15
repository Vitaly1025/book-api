package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"book-api/internal/domain"
	"book-api/internal/utils"
)

//go:generate mockgen -source=get_book_by_ids.go -destination=mocks/get_book_by_ids.go
type GetBookByIdsService interface {
	GetBookByIds(ctx context.Context, ids []int) ([]domain.BookEntity, error)
}

type GetBookByIdsHandler struct {
	bookService GetBookByIdsService
	logger      *slog.Logger
}

func NewGetBooksByIdsHandler(s GetBookByIdsService, l *slog.Logger) *GetBookByIdsHandler {
	return &GetBookByIdsHandler{bookService: s, logger: l}
}

// GetBookByIds godoc
// @Summary Retrieve books by their IDs
// @Description Retrieve a list of books by providing their IDs as a comma-separated query parameter.
// @Tags Books
// @Accept json
// @Produce json
// @Param ids query string true "Comma-separated list of book IDs"
// @Success 200 {array} dto.BookResponse "List of books"
// @Failure 400 {object} dto.ErrorResponse "Error retrieving books"
// @Router /books [get]
func (h *GetBookByIdsHandler) GetBookByIds(w http.ResponseWriter, r *http.Request) {
	const op = "GetBookById"
	log := h.logger.With(
		slog.String("operation", op),
	)

	ctx := context.Background()
	rw := utils.NewResponseWriter(w)

	query := r.URL.Query()
	idsParam := query.Get("ids")
	var ids []int
	for _, id := range strings.Split(idsParam, ",") {
		intId, _ := strconv.Atoi(id)
		ids = append(ids, intId)
	}

	resp, err := h.bookService.GetBookByIds(ctx, ids)
	if err != nil {
		log.Error("cannot get a book", slog.String("error", err.Error()))
		err = rw.Error(http.StatusBadRequest, "can't got a book")
		if err != nil {
			log.Error("cannot create error response", slog.String("error", err.Error()))
		}
	}

	err = rw.JSON(http.StatusOK, resp)
	if err != nil {
		log.Error("cannot create json response", slog.String("error", err.Error()))
	}
}
