package service_test

import (
	"context"
	"fmt"
	"testing"

	"book-api/internal/domain"
	mock_repository "book-api/internal/repository/mocks"
	"book-api/internal/service"
	"book-api/internal/transaction"
	mock_transaction "book-api/internal/transaction/mocks"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestBookService_CreateBook(t *testing.T) {
	type mockBehavior func(
		bookRepo *mock_repository.MockBookRepository,
		authorRepo *mock_repository.MockAuthorRepository,
		genreRepo *mock_repository.MockGenreRepository,
		txManager *mock_transaction.MockTransactionManager,
		ctx context.Context,
		book *domain.BookEntity,
	)

	testTable := []struct {
		name          string
		inputBook     domain.BookEntity
		mockBehavior  mockBehavior
		expectedID    int
		expectedError error
	}{
		{
			name: "OK",
			inputBook: domain.BookEntity{
				BookInfo: domain.BookInfo{
					Author: &domain.BookAuthor{Name: "Author Name"},
				},
				Genres: []domain.BookGenre{{Id: 1, Name: "Genre1"}},
			},
			mockBehavior: func(bookRepo *mock_repository.MockBookRepository, authorRepo *mock_repository.MockAuthorRepository, genreRepo *mock_repository.MockGenreRepository, txManager *mock_transaction.MockTransactionManager, ctx context.Context, book *domain.BookEntity) {
				authorRepo.EXPECT().GetByName(ctx, "Author Name").Return(nil, nil)
				authorRepo.EXPECT().CreateOrUpdate(ctx, gomock.Any(), book.BookInfo.Author).Return(nil)
				txManager.EXPECT().
					RunInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx transaction.Transaction) error) error {
						return fn(nil)
					})
				bookRepo.EXPECT().Create(ctx, gomock.Any(), book).Return(nil)
				book.Id = 1
				genreRepo.EXPECT().UpsertGenreToBook(ctx, gomock.Any(), book.Id, gomock.Any()).Return(nil)
			},
			expectedID:    1,
			expectedError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bookRepo := mock_repository.NewMockBookRepository(ctrl)
			authorRepo := mock_repository.NewMockAuthorRepository(ctrl)
			genreRepo := mock_repository.NewMockGenreRepository(ctrl)
			txManager := mock_transaction.NewMockTransactionManager(ctrl)

			testCase.mockBehavior(bookRepo, authorRepo, genreRepo, txManager, context.Background(), &testCase.inputBook)

			s := service.NewBookService(bookRepo, authorRepo, genreRepo, txManager)

			id, err := s.CreateBook(context.Background(), testCase.inputBook)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedID, id)
		})
	}
}

func TestBookService_UpdateBook(t *testing.T) {
	type mockBehavior func(
		bookRepo *mock_repository.MockBookRepository,
		authorRepo *mock_repository.MockAuthorRepository,
		genreRepo *mock_repository.MockGenreRepository,
		txManager *mock_transaction.MockTransactionManager,
		ctx context.Context,
		book domain.BookEntity,
	)

	testTable := []struct {
		name          string
		inputBook     domain.BookEntity
		mockBehavior  mockBehavior
		expectedID    int
		expectedError error
	}{
		{
			name: "OK",
			inputBook: domain.BookEntity{
				BookInfo: domain.BookInfo{
					Author: &domain.BookAuthor{Name: "Author Name"},
				},
				Genres: []domain.BookGenre{{Id: 1, Name: "Genre1"}},
			},
			mockBehavior: func(bookRepo *mock_repository.MockBookRepository, authorRepo *mock_repository.MockAuthorRepository, genreRepo *mock_repository.MockGenreRepository, txManager *mock_transaction.MockTransactionManager, ctx context.Context, book domain.BookEntity) {
				authorRepo.EXPECT().GetByName(ctx, "Author Name").Return(nil, nil)
				authorRepo.EXPECT().CreateOrUpdate(ctx, gomock.Any(), book.BookInfo.Author).Return(nil)
				txManager.EXPECT().
					RunInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx transaction.Transaction) error) error {
						return fn(nil)
					})
				bookRepo.EXPECT().Update(ctx, gomock.Any(), &book).Return(nil)
				genreRepo.EXPECT().UpsertGenreToBook(ctx, gomock.Any(), book.Id, gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bookRepo := mock_repository.NewMockBookRepository(ctrl)
			authorRepo := mock_repository.NewMockAuthorRepository(ctrl)
			genreRepo := mock_repository.NewMockGenreRepository(ctrl)
			txManager := mock_transaction.NewMockTransactionManager(ctrl)

			testCase.mockBehavior(bookRepo, authorRepo, genreRepo, txManager, context.Background(), testCase.inputBook)

			s := service.NewBookService(bookRepo, authorRepo, genreRepo, txManager)

			_, err := s.UpdateBook(context.Background(), testCase.inputBook)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestBookService_GetBookByIds(t *testing.T) {
	type mockBehavior func(
		bookRepo *mock_repository.MockBookRepository,
		authorRepo *mock_repository.MockAuthorRepository,
		genreRepo *mock_repository.MockGenreRepository,
		txManager *mock_transaction.MockTransactionManager,
		ctx context.Context,
		ids []int,
	)

	testTable := []struct {
		name          string
		inputIds      []int
		mockBehavior  mockBehavior
		expectedModel []domain.BookEntity
		expectedError error
	}{
		{
			name:     "OK",
			inputIds: []int{1, 2},
			mockBehavior: func(bookRepo *mock_repository.MockBookRepository, authorRepo *mock_repository.MockAuthorRepository, genreRepo *mock_repository.MockGenreRepository, txManager *mock_transaction.MockTransactionManager, ctx context.Context, ids []int) {
				books := []domain.BookEntity{
					{
						Id: 1,
						BookInfo: domain.BookInfo{
							Author:      &domain.BookAuthor{Name: "Author Name"},
							Description: "test",
							PageCount:   180,
							Rate:        9,
							Name:        "test1",
							Cover:       []byte("mock-image-data"),
						},
					},
					{
						Id: 2,
						BookInfo: domain.BookInfo{
							Author:      &domain.BookAuthor{Name: "Author Name"},
							Description: "test",
							PageCount:   180,
							Rate:        9,
							Name:        "test2",
							Cover:       []byte("mock-image-data"),
						},
					},
				}
				genres1 := []domain.BookGenre{
					{
						Name: "fantastic",
					},
					{
						Name: "horror",
					},
				}

				genres2 := []domain.BookGenre{
					{
						Name: "fantastic",
					},
				}
				bookRepo.EXPECT().GetByIds(ctx, ids).Return(books, nil)
				genreRepo.EXPECT().
					GetByBookId(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, bookId int) ([]domain.BookGenre, error) {
						if bookId == 1 {
							return genres1, nil
						}
						return genres2, nil
					}).
					Times(2)
			},
			expectedModel: getExpectedBookEntities(),
			expectedError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bookRepo := mock_repository.NewMockBookRepository(ctrl)
			authorRepo := mock_repository.NewMockAuthorRepository(ctrl)
			genreRepo := mock_repository.NewMockGenreRepository(ctrl)
			txManager := mock_transaction.NewMockTransactionManager(ctrl)

			testCase.mockBehavior(bookRepo, authorRepo, genreRepo, txManager, context.Background(), testCase.inputIds)

			s := service.NewBookService(bookRepo, authorRepo, genreRepo, txManager)

			actualResult, err := s.GetBookByIds(context.Background(), testCase.inputIds)
			assert.Equal(t, testCase.expectedError, err)
			for i := 0; i < len(testCase.expectedModel); i++ {
				expectedBook := testCase.expectedModel[i]
				actualBook := actualResult[i]

				assert.Equal(t, expectedBook, actualBook)
			}
		})
	}
}

func TestBookService_GetBooks(t *testing.T) {
	type mockBehavior func(
		bookRepo *mock_repository.MockBookRepository,
		authorRepo *mock_repository.MockAuthorRepository,
		genreRepo *mock_repository.MockGenreRepository,
		txManager *mock_transaction.MockTransactionManager,
		ctx context.Context,
	)

	testTable := []struct {
		name           string
		inputPredicate string
		mockBehavior   mockBehavior
		expectedModel  []domain.BookEntity
		expectedError  error
	}{
		{
			name:           "OK",
			inputPredicate: "tested",
			mockBehavior: func(bookRepo *mock_repository.MockBookRepository, authorRepo *mock_repository.MockAuthorRepository, genreRepo *mock_repository.MockGenreRepository, txManager *mock_transaction.MockTransactionManager, ctx context.Context) {
				books := []domain.BookEntity{
					{
						Id: 1,
						BookInfo: domain.BookInfo{
							Author:      &domain.BookAuthor{Name: "Author Name"},
							Description: "test",
							PageCount:   180,
							Rate:        9,
							Name:        "test1",
							Cover:       []byte("mock-image-data"),
						},
					},
					{
						Id: 2,
						BookInfo: domain.BookInfo{
							Author:      &domain.BookAuthor{Name: "Author Name"},
							Description: "test",
							PageCount:   180,
							Rate:        9,
							Name:        "test2",
							Cover:       []byte("mock-image-data"),
						},
					},
				}
				genres1 := []domain.BookGenre{
					{
						Name: "fantastic",
					},
					{
						Name: "horror",
					},
				}

				genres2 := []domain.BookGenre{
					{
						Name: "fantastic",
					},
				}

				bookRepo.EXPECT().GetAll(ctx).Return(books, nil)
				genreRepo.EXPECT().
					GetByBookId(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, bookId int) ([]domain.BookGenre, error) {
						if bookId == 1 {
							return genres1, nil
						}
						return genres2, nil
					}).
					Times(2)
			},
			expectedModel: getExpectedBookEntities(),
			expectedError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bookRepo := mock_repository.NewMockBookRepository(ctrl)
			authorRepo := mock_repository.NewMockAuthorRepository(ctrl)
			genreRepo := mock_repository.NewMockGenreRepository(ctrl)
			txManager := mock_transaction.NewMockTransactionManager(ctrl)

			testCase.mockBehavior(bookRepo, authorRepo, genreRepo, txManager, context.Background())

			s := service.NewBookService(bookRepo, authorRepo, genreRepo, txManager)

			actualResult, err := s.GetBooks(context.Background())
			assert.Equal(t, testCase.expectedError, err)
			for i := 0; i < len(testCase.expectedModel); i++ {
				expectedBook := testCase.expectedModel[i]
				actualBook := actualResult[i]

				assert.Equal(t, expectedBook, actualBook)
			}
		})
	}
}

func TestBookService_DeleteBook(t *testing.T) {
	type mockBehavior func(
		bookRepo *mock_repository.MockBookRepository,
		authorRepo *mock_repository.MockAuthorRepository,
		genreRepo *mock_repository.MockGenreRepository,
		txManager *mock_transaction.MockTransactionManager,
		ctx context.Context,
		id int,
	)

	testTable := []struct {
		name          string
		inputId       int
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name:    "OK",
			inputId: 1,
			mockBehavior: func(bookRepo *mock_repository.MockBookRepository, authorRepo *mock_repository.MockAuthorRepository, genreRepo *mock_repository.MockGenreRepository, txManager *mock_transaction.MockTransactionManager, ctx context.Context, id int) {
				txManager.EXPECT().
					RunInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx transaction.Transaction) error) error {
						return fn(nil)
					})
				bookRepo.EXPECT().Delete(ctx, gomock.Any(), id).Return(nil)
				genreRepo.EXPECT().DeleteBookGenre(ctx, gomock.Any(), id).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:    "Delete Book Error",
			inputId: 1,
			mockBehavior: func(bookRepo *mock_repository.MockBookRepository, authorRepo *mock_repository.MockAuthorRepository, genreRepo *mock_repository.MockGenreRepository, txManager *mock_transaction.MockTransactionManager, ctx context.Context, id int) {
				txManager.EXPECT().
					RunInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx transaction.Transaction) error) error {
						return fn(nil)
					})
				genreRepo.EXPECT().DeleteBookGenre(ctx, gomock.Any(), id).Return(nil)
				bookRepo.EXPECT().Delete(ctx, gomock.Any(), id).Return(fmt.Errorf("error deleting book"))
			},
			expectedError: fmt.Errorf("failed to remove book: %s", "error deleting book"),
		},
		{
			name:    "Delete Book Genre Error",
			inputId: 1,
			mockBehavior: func(bookRepo *mock_repository.MockBookRepository, authorRepo *mock_repository.MockAuthorRepository, genreRepo *mock_repository.MockGenreRepository, txManager *mock_transaction.MockTransactionManager, ctx context.Context, id int) {
				txManager.EXPECT().
					RunInTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(tx transaction.Transaction) error) error {
						return fn(nil)
					})
				genreRepo.EXPECT().DeleteBookGenre(ctx, gomock.Any(), id).Return(fmt.Errorf("error deleting book genre"))
			},
			expectedError: fmt.Errorf("failed to remove bookgenre: %s", "error deleting book genre"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bookRepo := mock_repository.NewMockBookRepository(ctrl)
			authorRepo := mock_repository.NewMockAuthorRepository(ctrl)
			genreRepo := mock_repository.NewMockGenreRepository(ctrl)
			txManager := mock_transaction.NewMockTransactionManager(ctrl)

			ctx := context.Background()
			testCase.mockBehavior(bookRepo, authorRepo, genreRepo, txManager, ctx, testCase.inputId)

			s := service.NewBookService(bookRepo, authorRepo, genreRepo, txManager)

			err := s.DeleteBook(ctx, testCase.inputId)
			if testCase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func getExpectedBookEntities() []domain.BookEntity {
	return []domain.BookEntity{
		{
			Id: 1,
			BookInfo: domain.BookInfo{
				Author:      &domain.BookAuthor{Name: "Author Name"},
				Description: "test",
				PageCount:   180,
				Rate:        9,
				Name:        "test1",
				Cover:       []byte("mock-image-data"),
			},
			Genres: []domain.BookGenre{
				{
					Name: "fantastic",
				},
				{
					Name: "horror",
				},
			},
		},
		{
			Id: 2,
			BookInfo: domain.BookInfo{
				Author:      &domain.BookAuthor{Name: "Author Name"},
				Description: "test",
				PageCount:   180,
				Rate:        9,
				Name:        "test2",
				Cover:       []byte("mock-image-data"),
			},
			Genres: []domain.BookGenre{
				{
					Name: "fantastic",
				},
			},
		},
	}
}
