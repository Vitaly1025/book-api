package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"book-api/internal/config"
	handlers "book-api/internal/handler/book"
	"book-api/internal/middleware"
	"book-api/internal/repository/sqlx"
	"book-api/internal/service"
	"book-api/internal/transaction"

	_ "book-api/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Book API
// @version 1.0
// @description This is a sample server.
// @host localhost:4000
func main() {
	cfg := config.MustLoad()
	log := setupLogger(slog.LevelDebug)

	log.Debug("debug messages are enabled")

	db, err := sqlx.NewPostgresDB(cfg.Db)
	if err != nil {
		log.Error("failed to establish a connection to db", slog.Any("err", err))
	}

	middlewares := []middleware.Middleware{
		middleware.NewLoggerMiddlware(log),
	}

	bookRepo := sqlx.NewSQLXBookRepository(db)
	authorRepo := sqlx.NewSQLXAuthorRepository(db)
	bookGenreRepo := sqlx.NewSQLXAGenreRepository(db)
	txManager := transaction.NewTransactionManager(db)

	bookService := service.NewBookService(bookRepo, authorRepo, bookGenreRepo, txManager)
	createBookHandler := handlers.NewCreateBookHandler(bookService, log)
	updateBookHandler := handlers.NewUpdateBookHandler(bookService, log)
	deleteBookHandler := handlers.NewDeleteBookHandler(bookService, log)
	getAllBookHandler := handlers.NewGetAllBookHandler(bookService, log)
	getByIdsBookHandler := handlers.NewGetBooksByIdsHandler(bookService, log)

	router := mux.NewRouter()

	// swagger initialization
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// init basic routes
	router.HandleFunc("/book", createBookHandler.CreateBook).Methods(http.MethodPost)
	router.HandleFunc("/book", updateBookHandler.UpdateBook).Methods(http.MethodPut)
	router.HandleFunc("/book/{id}", getByIdsBookHandler.GetBookByIds).Methods(http.MethodGet)
	router.HandleFunc("/book/}", getAllBookHandler.GetAllBook).Methods(http.MethodGet)
	router.HandleFunc("/book/{id}", deleteBookHandler.DeleteBook).Methods(http.MethodDelete)

	// init middlewares
	for _, middleware := range middlewares {
		router.Use(middleware.Proccess)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", cfg.HttpServer.Host, cfg.HttpServer.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", slog.Any("err", err))
		}
	}()

	log.Info("server started")
	<-done
	log.Info("stopping server")

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.Any("err", err))

		return
	}

	// TODO: close storage
	log.Info("server stopped")
}

func setupLogger(lvl slog.Leveler) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
}
