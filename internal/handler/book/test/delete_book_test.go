package handlers_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	handlers "book-api/internal/handler/book"
	mock_handlers "book-api/internal/handler/book/mocks"

	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
	"go.uber.org/mock/gomock"
)

func Test_DeleteBookHandler(t *testing.T) {
	type mockBehavior func(s *mock_handlers.MockDeleteBookService, ctx context.Context, id int)

	testTable := []struct {
		name                string
		inputParam          string
		inputRequest        int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:         "OK",
			inputParam:   "1",
			inputRequest: 1,
			mockBehavior: func(s *mock_handlers.MockDeleteBookService, ctx context.Context, id int) {
				s.EXPECT().DeleteBook(ctx, id).Return(nil)
			},
			expectedStatusCode: 204,
		},
		// {
		// 	name:                "Incorrcet Param",
		// 	inputParam:          "sfs",
		// 	inputRequest:        1,
		// 	mockBehavior:        func(s *mock_handlers.MockDeleteBookService, id int) {},
		// 	expectedStatusCode:  400,
		// 	expectedRequestBody: `{"message":"incorrect param"}` + "\n",
		// },
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			bookService := mock_handlers.NewMockDeleteBookService(c)
			testCase.mockBehavior(bookService, context.Background(), testCase.inputRequest)

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
			handler := handlers.NewDeleteBookHandler(bookService, logger)

			r := mux.NewRouter()
			r.HandleFunc("/book/{id}", handler.DeleteBook).Methods(http.MethodDelete)

			w := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/book/%s", testCase.inputParam), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedRequestBody)
		})
	}
}
