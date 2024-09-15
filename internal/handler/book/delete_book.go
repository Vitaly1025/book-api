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

// DeleteBook godoc
// @Summary Delete a book by ID
// @Description Delete a book entry from the system using the book's ID
// @Tags Books
// @Accept json
// @Produce json
// @Param id query int true "Book ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse "Bad Request: Invalid parameter or deletion issue"
// @Router /books [delete]
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
		err = rw.Error(http.StatusBadRequest, "can't parse param")
		if err != nil {
			log.Error("cannot create error response", slog.String("error", err.Error()))
		}
		return
	}

	intId, err := strconv.Atoi(id)
	if err != nil {
		log.Error("cannot parse to int", slog.String("error", err.Error()))
		err = rw.Error(http.StatusBadRequest, "incorrect param")
		if err != nil {
			log.Error("cannot create error response", slog.String("error", err.Error()))
		}
		return
	}

	err = h.bookService.DeleteBook(ctx, intId)
	if err != nil {
		log.Error("cannot delete a book", slog.String("error", err.Error()))
		err = rw.Error(http.StatusBadRequest, "problem with deleting the book")
		if err != nil {
			log.Error("cannot create error response", slog.String("error", err.Error()))
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
