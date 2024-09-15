package utils

import (
	dto "book-api/pkg/dto/http"
	"encoding/json"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w}
}

func (w *ResponseWriter) WriteStatus(statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func (w *ResponseWriter) JSON(statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func (w *ResponseWriter) Error(statusCode int, message string) error {
	return w.JSON(statusCode, dto.ErrorResponse{Error: message})
}
