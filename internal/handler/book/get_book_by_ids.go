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

// @Summary Get book by ID
// @Tags Book Operations
// @Description This method gets book via id
// @Accept text/json
// @Param id path int true "Book Ids"
// @Produce  json
// @Router /book/{id} [get]
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
		rw.Error(http.StatusBadRequest, "can't got a book")
	}

	rw.JSON(http.StatusOK, resp)
}
