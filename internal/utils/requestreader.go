package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type RequestReader struct{}

func NewRequestReader() *RequestReader {
	return &RequestReader{}
}

func (r *RequestReader) ParseGetRequest(req *http.Request, model interface{}) error {
	queryParams := req.URL.Query()
	jsonBytes, err := json.Marshal(queryParams)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, model)
}

func (r *RequestReader) ParseDeleteRequest(req *http.Request, model interface{}) error {
	vars := mux.Vars(req)

	jsonBytes, err := json.Marshal(vars)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonBytes, model)
}

func (r *RequestReader) ParsePostRequest(req *http.Request, model interface{}) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	return json.Unmarshal(body, model)
}
