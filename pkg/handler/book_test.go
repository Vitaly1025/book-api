package handlers

import (
	"book-api/pkg/middleware"
	models "book-api/pkg/model"
	"book-api/pkg/service"
	mock_service "book-api/pkg/service/mocks"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
)

func Test_CreateBookHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookWork, req models.BookRequest)

	testTable := []struct {
		name                string
		inputBody           string
		inputRequest        models.BookRequest
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name": "Test", "genre": 1, "price": 20, "amount": 10}`,
			inputRequest: models.BookRequest{
				Name:   "Test",
				Genre:  1,
				Price:  20,
				Amount: 10,
			},
			mockBehavior: func(s *mock_service.MockBookWork, book models.BookRequest) {
				s.EXPECT().CreateBook(book).Return(&models.SimpleResponse{Id: 1}, nil)
			},
			expectedStatusCode: 200,
			//This is because Encoder add extra new line in the end
			//https://stackoverflow.com/questions/36319918/why-does-json-encoder-add-an-extra-line/36320146
			expectedRequestBody: `{"id":1}` + "\n",
		},
		{
			name:                "Invalid validation",
			inputBody:           `{"name": "Test2", "genre": 1, "price": -20, "amount": -10}`,
			mockBehavior:        func(s *mock_service.MockBookWork, book models.BookRequest) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"data validation failed"}` + "\n",
		},
		{
			name:                "Incorrect data",
			inputBody:           `{"id": 10}`,
			mockBehavior:        func(s *mock_service.MockBookWork, book models.BookRequest) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid data"}` + "\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			bookService := mock_service.NewMockBookWork(c)
			testCase.mockBehavior(bookService, testCase.inputRequest)

			mdw := middleware.NewMiddleware()
			services := &service.Service{BookWork: bookService}
			handler := NewHandler(services, mdw)

			r := mux.NewRouter()
			r.HandleFunc("/create-book", handler.CreateBook).Methods(http.MethodPost)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/create-book", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
func Test_UpdateBookHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookWork, req models.Book)

	testTable := []struct {
		name                string
		inputBody           string
		inputRequest        models.Book
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"id": 1,"name": "Test", "genre": 1, "price": 20, "amount": 10}`,
			inputRequest: models.Book{
				Id: 1,
				BookRequest: models.BookRequest{
					Name:   "Test",
					Genre:  1,
					Price:  20,
					Amount: 10,
				},
			},
			mockBehavior: func(s *mock_service.MockBookWork, book models.Book) {
				s.EXPECT().UpdateBook(book).Return(&models.SimpleResponse{Id: 1}, nil)
			},
			expectedStatusCode: 200,
			//This is because Encoder add extra new line in the end
			//https://stackoverflow.com/questions/36319918/why-does-json-encoder-add-an-extra-line/36320146
			expectedRequestBody: `{"id":1}` + "\n",
		},
		{
			name:                "Invalid validation",
			inputBody:           `{"name": "Test2", "genre": 1, "price": -20, "amount": -10}`,
			mockBehavior:        func(s *mock_service.MockBookWork, book models.Book) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"data validation failed"}` + "\n",
		},
		{
			name:                "Incorrect data",
			inputBody:           `{"tee": 10}`,
			mockBehavior:        func(s *mock_service.MockBookWork, book models.Book) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid data"}` + "\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			bookService := mock_service.NewMockBookWork(c)
			testCase.mockBehavior(bookService, testCase.inputRequest)

			mdw := middleware.NewMiddleware()
			services := &service.Service{BookWork: bookService}
			handler := NewHandler(services, mdw)

			r := mux.NewRouter()
			r.HandleFunc("/update-book", handler.UpdateBook).Methods(http.MethodPost)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/update-book", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func Test_DeleteBookHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookWork, id int)

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
			mockBehavior: func(s *mock_service.MockBookWork, id int) {
				s.EXPECT().DeleteBook(id).Return(nil)
			},
			expectedStatusCode: 204,
		},
		{
			name:                "Incorrcet Param",
			inputParam:          "sfs",
			inputRequest:        1,
			mockBehavior:        func(s *mock_service.MockBookWork, id int) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"incorrect param"}` + "\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			bookService := mock_service.NewMockBookWork(c)
			testCase.mockBehavior(bookService, testCase.inputRequest)

			mdw := middleware.NewMiddleware()
			services := &service.Service{BookWork: bookService}
			handler := NewHandler(services, mdw)

			r := mux.NewRouter()
			r.HandleFunc("/delete-book/{id}", handler.DeleteBook).Methods(http.MethodDelete)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/delete-book/%s", testCase.inputParam), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func Test_GetBookByIdHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookWork, id int)

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
			mockBehavior: func(s *mock_service.MockBookWork, id int) {
				s.EXPECT().GetBookById(id).Return(&models.Book{}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":0,"name":"","price":0,"genre":0,"amount":0}` + "\n",
		},
		{
			name:                "Incorrcet Param",
			inputParam:          "sfs",
			inputRequest:        1,
			mockBehavior:        func(s *mock_service.MockBookWork, id int) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"incorrect param"}` + "\n",
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			bookService := mock_service.NewMockBookWork(c)
			testCase.mockBehavior(bookService, testCase.inputRequest)

			mdw := middleware.NewMiddleware()
			services := &service.Service{BookWork: bookService}
			handler := NewHandler(services, mdw)

			r := mux.NewRouter()
			r.HandleFunc("/get-book/{id}", handler.GetBookById).Methods(http.MethodGet)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/get-book/%s", testCase.inputParam), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func Test_GetAllBookHandler(t *testing.T) {
	type mockBehavior func(s *mock_service.MockBookWork, bookname string, genre string)

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
				"genre": "", "bookname": "",
			},
			inputRequest: 1,
			mockBehavior: func(s *mock_service.MockBookWork, bookname string, genre string) {
				s.EXPECT().GetBooks(bookname, genre).Return(&[]models.Book{}, nil)
			},
			expectedStatusCode: 200,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			bookService := mock_service.NewMockBookWork(c)
			testCase.mockBehavior(bookService, testCase.params["bookname"], testCase.params["genre"])

			mdw := middleware.NewMiddleware()
			services := &service.Service{BookWork: bookService}
			handler := NewHandler(services, mdw)

			r := mux.NewRouter()
			r.HandleFunc("/get-books", handler.GetAllBook).Methods(http.MethodGet)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/get-books", nil)

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
