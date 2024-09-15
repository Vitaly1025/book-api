package middleware

import "net/http"

type Middleware interface {
	Proccess(next http.Handler) http.Handler
}
