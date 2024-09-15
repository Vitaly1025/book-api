package handlers_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"book-api/internal/domain"
	handlers "book-api/internal/handler/book"
	mock_handlers "book-api/internal/handler/book/mocks"

	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"
)

func Test_GetAllBookHandler(t *testing.T) {
	type mockBehavior func(s *mock_handlers.MockGetAllBookService, ctx context.Context)

	testTable := []struct {
		name                string
		params              map[string]string
		inputRequest        int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			params: map[string]string{
				"predicate": "",
			},
			inputRequest: 1,
			mockBehavior: func(s *mock_handlers.MockGetAllBookService, ctx context.Context) {
				s.EXPECT().GetBooks(ctx).Return([]domain.BookEntity{}, nil)
			},
			expectedStatusCode: 200,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			bookService := mock_handlers.NewMockGetAllBookService(c)
			testCase.mockBehavior(bookService, context.Background())

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
			handler := handlers.NewGetAllBookHandler(bookService, logger)

			r := mux.NewRouter()
			r.HandleFunc("/books", handler.GetAllBook).Methods(http.MethodGet)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/books", nil)

			q := req.URL.Query()
			for k, v := range testCase.params {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			// assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
