package handlers_test

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"book-api/internal/domain"
	handlers "book-api/internal/handler/book"
	mock_handlers "book-api/internal/handler/book/mocks"
	"book-api/internal/mapper"
	dto "book-api/pkg/dto/http"

	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"
)

func Test_CreateBookHandler(t *testing.T) {
	type mockBehavior func(s *mock_handlers.MockCreateBookService, ctx context.Context, req domain.BookEntity)
	handlerPath := "/book"

	testTable := []struct {
		name                string
		inputBody           string
		inputRequest        dto.CreateBookRequest
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			inputBody: `{
			"name": "Sherlock Holmes", 
			"genres": [1,2], 
			"author_name": "Arthur Conan Doyle", 
			"description": "This book present a range of cases, from murder mysteries to thefts and conspiracies, showcasing Holmes’ extraordinary abilities of observation, deduction, and analytical thinking.",
			"rate": 9,
			"cover": "bW9jay1pbWFnZS1kYXRh",
			"page_count": 182}`,
			inputRequest: dto.CreateBookRequest{
				Name:        "Sherlock Holmes",
				Genres:      []int{1, 2},
				AuthorName:  "Arthur Conan Doyle",
				Description: "This book present a range of cases, from murder mysteries to thefts and conspiracies, showcasing Holmes’ extraordinary abilities of observation, deduction, and analytical thinking.",
				Cover:       []byte("mock-image-data"),
				Rate:        9,
				PageCount:   182,
			},
			mockBehavior: func(s *mock_handlers.MockCreateBookService, ctx context.Context, book domain.BookEntity) {
				s.EXPECT().CreateBook(ctx, book).Return(1, nil)
			},
			expectedStatusCode: 200,
			// This is because Encoder add extra new line in the end
			// https://stackoverflow.com/questions/36319918/why-does-json-encoder-add-an-extra-line/36320146
			expectedRequestBody: `{"id":1}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			bookService := mock_handlers.NewMockCreateBookService(c)
			testCase.mockBehavior(bookService, context.Background(), mapper.ToDomainCreateBook(testCase.inputRequest))

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
			handler := handlers.NewCreateBookHandler(bookService, logger)

			r := mux.NewRouter()
			r.HandleFunc(handlerPath, handler.CreateBook).Methods(http.MethodPost)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, handlerPath, bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedRequestBody)
		})
	}
}
