package middleware

import (
	models "book-api/model"
	"encoding/json"
	"log"
	"net/http"
)

type ErrorMiddleware struct {
}

func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

func (m *ErrorMiddleware) Error(msg string, systemMessage string, sts int, w http.ResponseWriter){
	err := models.Error {
		Message: msg,
	}
	log.Println(systemMessage)
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(err)
}