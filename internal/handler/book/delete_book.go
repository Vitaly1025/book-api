package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"book-api/internal/utils"
)

//go:generate mockgen -source=delete_book.go -destination=mocks/delete_book.go
type DeleteBookService interface {
	DeleteBook(ctx context.Context, id int) error
}

type DeleteBookHandler struct {
	bookService DeleteBookService
	logger      *slog.Logger
}

func NewDeleteBookHandler(s DeleteBookService, l *slog.Logger) *DeleteBookHandler {
	return &DeleteBookHandler{bookService: s, logger: l}
}

// @Summary Delete book by ID
// @Tags Book Operations
// @Description This method delete book by id
// @Accept text/json
// @Param id path int true "Book Id"
// @Produce  json
// @Router /book/{id} [delete]
// @Success  204  {int} resp
func (h *DeleteBookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	const op = "DeleteBook"
	log := h.logger.With(
		slog.String("operation", op),
	)
	ctx := context.Background()
	rw := utils.NewResponseWriter(w)
	rr := utils.NewRequestReader()

	req := make(map[string]string)
	err := rr.ParseDeleteRequest(r, &req)

	id, ok := req["id"]
	if !ok {
		log.Error("cannot find 'id' from request", slog.String("error", err.Error()))
		rw.Error(http.StatusBadRequest, "can't parse param")
		return
	}

	intId, err := strconv.Atoi(id)
	if err != nil {
		log.Error("cannot parse to int", slog.String("error", err.Error()))
		rw.Error(http.StatusBadRequest, "incorrect param")
		return
	}

	err = h.bookService.DeleteBook(ctx, intId)
	if err != nil {
		log.Error("cannot delete a book", slog.String("error", err.Error()))
		rw.Error(http.StatusBadRequest, "problem with deleting the book")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
