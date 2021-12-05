package handlers

import (
	models "book-api/pkg/model"
	validators "book-api/pkg/validator"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// @Summary Create a book
// @Description Create a book with data
// @Accept text/json
// @Param request body models.BookRequest true "Book"
// @Produce  json
// @Router /create-book [post]
// @Success  200  {int} resp
func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p models.BookRequest

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.md.Error.Error("invalid data", err.Error(), http.StatusBadRequest, w)
		return
	} else if(models.BookRequest{}) == p {
		h.md.Error.Error("invalid data", "", http.StatusBadRequest, w)
		return
	}

	err := validators.ValidateRequest(p)
	if err != nil {
		h.md.Error.Error("data validation failed", err.Error(), http.StatusBadRequest, w)
		return
	}

	resp, err := h.services.BookWork.CreateBook(p)
	if err != nil {
		h.md.Error.Error("smth went wrong :(", err.Error(), http.StatusBadRequest, w)
		return
	}
	log.Println(resp)
	json.NewEncoder(w).Encode(resp)
}

// @Summary Get book by ID
// @Description This method gets book via id
// @Accept text/json
// @Param id path int true "Book Id"
// @Produce  json
// @Router /get-book/{id} [get]
func (h *Handler) GetBookById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    id, ok := vars["id"]
	if !ok {
		h.md.Error.Error("Can't parse param", "", http.StatusBadRequest, w)
		return
	}
	
	intId, err := strconv.Atoi(id)
	if err != nil{
		h.md.Error.Error("incorrect param", err.Error(), http.StatusBadRequest, w)
		return
	}
	
	resp, err := h.services.BookWork.GetBookById(intId)
	if err != nil {
		h.md.Error.Error("can't got a book", err.Error(), http.StatusBadRequest, w)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

// @Summary Delete book by ID
// @Description This method delete book by id
// @Accept text/json
// @Param id path int true "Book Id"
// @Produce  json
// @Router /delete-book/{id} [delete]
// @Success  204  {int} resp
func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    id, ok := vars["id"]
	if !ok {
		h.md.Error.Error("Can't parse param", "", http.StatusBadRequest, w)
		return
	}
	
    intId, err := strconv.Atoi(id)
	if err != nil{
		h.md.Error.Error("incorrect param", err.Error(), http.StatusBadRequest, w)
		return
	}
	err = h.services.BookWork.DeleteBook(intId)
	if err != nil {
		h.md.Error.Error("problem with deleting book", err.Error(), http.StatusBadRequest, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Get books
// @Description This method return books
// @Accept text/json
// @Param bookname query string false "BookName"
// @Param genre query string false "GenreName"
// @Produce  json
// @Router /get-books [get]
func (h *Handler) GetAllBook(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
    bookname := query.Get("bookname")
	genre := query.Get("genre")

	resp, err := h.services.BookWork.GetBooks(bookname, genre)
	if err != nil {
		h.md.Error.Error("can't got a book", err.Error(), http.StatusBadRequest, w)
		return
	}
	json.NewEncoder(w).Encode(resp)
}


// @Summary Update book
// @Tags Book Operations
// @Description Update book
// @Accept text/json
// @Param request body models.Book true "Book"
// @Produce  json
// @Router /update-book [post]
// @Success  200  {int} resp
func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var p models.Book

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		h.md.Error.Error("invalid data", err.Error(), http.StatusBadRequest, w)
		return
	} else if(models.Book{}) == p {
		h.md.Error.Error("invalid data", "", http.StatusBadRequest, w)
		return
	}

	err := validators.ValidateRequest(p.BookRequest)
	if err != nil {
		h.md.Error.Error("data validation failed", err.Error(), http.StatusBadRequest, w)
		return
	}

	resp, err := h.services.BookWork.UpdateBook(p)
	if err != nil {
		h.md.Error.Error("smth went wrong :(",err.Error(), http.StatusBadRequest, w)
		return
	}

	json.NewEncoder(w).Encode(resp)
}
