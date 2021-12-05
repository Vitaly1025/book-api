package handlers

import (
	"net/http"

	_ "book-api/docs"
	"book-api/pkg/middleware"
	"book-api/pkg/service"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct{
	services *service.Service
	md *middleware.Middleware
}

func NewHandler(s *service.Service, md *middleware.Middleware) *Handler{
	return &Handler{services: s, md: md}
}

func (h *Handler) InitRoutes() *mux.Router{
	router := mux.NewRouter()

	//Swagger Initialization
	router.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	//Init basic routes
	router.HandleFunc("/create-book", h.CreateBook).Methods(http.MethodPost)
	router.HandleFunc("/update-book", h.UpdateBook).Methods(http.MethodPost)
	router.HandleFunc("/get-book/{id}", h.GetBookById).Methods(http.MethodGet)
	router.HandleFunc("/get-books", h.GetAllBook).Methods(http.MethodGet)
	router.HandleFunc("/delete-book/{id}", h.DeleteBook).Methods(http.MethodDelete)
	

	return router
}