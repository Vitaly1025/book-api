package middleware

import (
	"net/http"
)

type Error interface {
	Error(msg string, systemMessage string, sts int, w http.ResponseWriter)
}

type ValidationMiddleware interface {
	ParseRequest()
}

type Middleware struct {
	Error
}

func NewMiddleware() *Middleware { 
	return &Middleware{
		Error: NewErrorMiddleware(),
	} 
}