package sqlx_test

import (
	"context"
	"fmt"
	"testing"

	"book-api/internal/domain"
	"book-api/internal/repository/sqlx"
	"book-api/internal/transaction"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	provider "github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := provider.NewDb(db, "sqlmock")
	txManager := transaction.NewTransactionManager(sqlxDB)
	repo := sqlx.NewSQLXBookRepository(sqlxDB)
	ctx := context.Background()

	assert.NoError(t, err)

	tests := []struct {
		testName       string
		book           *domain.BookEntity
		mockSetup      func()
		expectErr      bool
		expectedBookID int
	}{
		{
			testName: "Successful book creation",
			book: &domain.BookEntity{
				BookInfo: domain.BookInfo{
					Author:      &domain.BookAuthor{Id: 1},
					Name:        "Test Book",
					Description: "Test Description",
					Cover:       []byte("cover.png"),
					PageCount:   100,
					Rate:        4.5,
				},
			},
			mockSetup: func() {
				query := fmt.Sprintf(`INSERT INTO %s \(name, description, cover, page_count, rate, author_id\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING id`, sqlx.BookTable)
				mock.ExpectBegin()
				mock.ExpectQuery(query).
					WithArgs("Test Book", "Test Description", []byte("cover.png"), 100, 4.5, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectedBookID: 1,
			expectErr:      false,
		},
		{
			testName: "Error during book creation",
			book: &domain.BookEntity{
				BookInfo: domain.BookInfo{
					Author:      &domain.BookAuthor{Id: 1},
					Name:        "Test Book",
					Description: "Test Description",
					Cover:       []byte("cover.png"),
					PageCount:   100,
					Rate:        4.5,
				},
			},
			mockSetup: func() {
				query := fmt.Sprintf(`INSERT INTO %s \(name, description, cover, page_count, rate, author_id\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING id`, sqlx.BookTable)
				expectedErr := fmt.Errorf("db error")
				mock.ExpectBegin()
				mock.ExpectQuery(query).
					WithArgs("Test Book", "Test Description", []byte("cover.png"), 100, 4.5, 1).
					WillReturnError(expectedErr)
				mock.ExpectRollback().WillReturnError(expectedErr)
			},
			expectedBookID: 0,
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.mockSetup()

			err := txManager.RunInTransaction(ctx, func(tx transaction.Transaction) error {
				return repo.Create(ctx, tx, tt.book)
			})

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBookID, tt.book.Id)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := provider.NewDb(db, "sqlmock")
	txManager := transaction.NewTransactionManager(sqlxDB)
	repo := sqlx.NewSQLXBookRepository(sqlxDB)
	ctx := context.Background()

	assert.NoError(t, err)

	tests := []struct {
		testName       string
		book           *domain.BookEntity
		mockSetup      func()
		expectErr      bool
		expectedBookID int
	}{
		{
			testName: "Successful book creation",
			book: &domain.BookEntity{
				BookInfo: domain.BookInfo{
					Author:      &domain.BookAuthor{Id: 1},
					Name:        "Test Book",
					Description: "Test Description",
					Cover:       []byte("cover.png"),
					PageCount:   100,
					Rate:        4.5,
				},
			},
			mockSetup: func() {
				query := fmt.Sprintf(`UPDATE %s SET description = \$2, cover = \$3, page_count = \$4, rate = \$5, author_id = \$6
				WHERE name = \$1
				RETURNING id`, sqlx.BookTable)
				mock.ExpectBegin()
				mock.ExpectQuery(query).
					WithArgs("Test Book", "Test Description", []byte("cover.png"), 100, 4.5, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectedBookID: 1,
			expectErr:      false,
		},
		{
			testName: "Error during book updating",
			book: &domain.BookEntity{
				BookInfo: domain.BookInfo{
					Author:      &domain.BookAuthor{Id: 1},
					Name:        "Test Book",
					Description: "Test Description",
					Cover:       []byte("cover.png"),
					PageCount:   100,
					Rate:        4.5,
				},
			},
			mockSetup: func() {
				query := fmt.Sprintf(`UPDATE %s SET description = \$2, cover = \$3, page_count = \$4, rate = \$5, author_id = \$6
				WHERE name = \$1
				RETURNING id`, sqlx.BookTable)
				expectedErr := fmt.Errorf("db error")
				mock.ExpectBegin()
				mock.ExpectQuery(query).
					WithArgs("Test Book", "Test Description", []byte("cover.png"), 100, 4.5, 1).
					WillReturnError(expectedErr)
				mock.ExpectRollback().WillReturnError(expectedErr)
			},
			expectedBookID: 0,
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.mockSetup()

			err := txManager.RunInTransaction(ctx, func(tx transaction.Transaction) error {
				return repo.Update(ctx, tx, tt.book)
			})

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBookID, tt.book.Id)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetByIdsBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := provider.NewDb(db, "sqlmock")
	repo := sqlx.NewSQLXBookRepository(sqlxDB)
	ctx := context.Background()

	assert.NoError(t, err)

	tests := []struct {
		name      string
		mockSetup func()
		bookIds   []int
		expected  []domain.BookEntity
		expectErr bool
	}{
		{
			name:    "Successful retrieval by IDs",
			bookIds: []int{1, 2},
			mockSetup: func() {
				query := fmt.Sprintf(`
				SELECT 
					b.id as book_id, 
					b.cover, 
					b.description, 
					b.name as book_name, 
					b.page_count, 
					b.rate, 
					a.id as author_id, 
					a.name as author_name
				FROM %s AS b
				LEFT JOIN author AS a ON b.author_id = a.id
				WHERE b.id IN (?)`, sqlx.BookTable)
				mock.ExpectQuery(query).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"book_id", "cover", "description", "book_name", "page_count", "rate", "author_id", "author_name"}).
						AddRow(1, "cover1.png", "Description 1", "Book 1", 100, 4.5, 1, "Author 1").
						AddRow(2, "cover2.png", "Description 2", "Book 2", 200, 4.7, 2, "Author 2"))
			},
			expected: []domain.BookEntity{
				{
					Id: 1,
					BookInfo: domain.BookInfo{
						Cover:       []byte("cover1.png"),
						Description: "Description 1",
						Name:        "Book 1",
						PageCount:   100,
						Rate:        4.5,
						Author: &domain.BookAuthor{
							Id:   1,
							Name: "Author 1",
						},
					},
				},
				{
					Id: 2,
					BookInfo: domain.BookInfo{
						Cover:       []byte("cover2.png"),
						Description: "Description 2",
						Name:        "Book 2",
						PageCount:   200,
						Rate:        4.7,
						Author: &domain.BookAuthor{
							Id:   2,
							Name: "Author 2",
						},
					},
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			books, err := repo.GetByIds(ctx, tt.bookIds)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.expected, books)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := provider.NewDb(db, "sqlmock")
	txManager := transaction.NewTransactionManager(sqlxDB)
	repo := sqlx.NewSQLXBookRepository(sqlxDB)
	ctx := context.Background()

	assert.NoError(t, err)

	tests := []struct {
		testName  string
		bookId    int
		mockSetup func()
		expectErr bool
	}{
		{
			testName: "Successful book deletion",
			bookId:   1,
			mockSetup: func() {
				query := fmt.Sprintf(`DELETE FROM %s WHERE id = \$1`, sqlx.BookTable)
				mock.ExpectBegin()
				mock.ExpectExec(query).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			testName: "Error during book deletion",
			bookId:   1,
			mockSetup: func() {
				query := fmt.Sprintf(`DELETE FROM %s WHERE id = \$1`, sqlx.BookTable)
				expectedErr := fmt.Errorf("db error")
				mock.ExpectBegin()
				mock.ExpectExec(query).WithArgs(1).WillReturnError(expectedErr)
				mock.ExpectRollback().WillReturnError(expectedErr)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.mockSetup()

			err := txManager.RunInTransaction(ctx, func(tx transaction.Transaction) error {
				return repo.Delete(ctx, tx, tt.bookId)
			})

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAllBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := provider.NewDb(db, "sqlmock")
	repo := sqlx.NewSQLXBookRepository(sqlxDB)
	ctx := context.Background()

	assert.NoError(t, err)

	tests := []struct {
		name      string
		mockSetup func()
		expected  []domain.BookEntity
		expectErr bool
	}{
		{
			name: "Successful retrieval by IDs",
			mockSetup: func() {
				query := fmt.Sprintf(`
				SELECT 
					b.id as book_id, 
					b.cover, 
					b.description, 
					b.name as book_name, 
					b.page_count, 
					b.rate, 
					a.id as author_id, 
					a.name as author_name
				FROM %s AS b
				LEFT JOIN author AS a ON b.author_id = a.id;`, sqlx.BookTable)
				mock.ExpectQuery(query).
					WillReturnRows(sqlmock.NewRows([]string{"book_id", "cover", "description", "book_name", "page_count", "rate", "author_id", "author_name"}).
						AddRow(1, "cover1.png", "Description 1", "Book 1", 100, 4.5, 1, "Author 1").
						AddRow(2, "cover2.png", "Description 2", "Book 2", 200, 4.7, 2, "Author 2"))
			},
			expected: []domain.BookEntity{
				{
					Id: 1,
					BookInfo: domain.BookInfo{
						Cover:       []byte("cover1.png"),
						Description: "Description 1",
						Name:        "Book 1",
						PageCount:   100,
						Rate:        4.5,
						Author: &domain.BookAuthor{
							Id:   1,
							Name: "Author 1",
						},
					},
				},
				{
					Id: 2,
					BookInfo: domain.BookInfo{
						Cover:       []byte("cover2.png"),
						Description: "Description 2",
						Name:        "Book 2",
						PageCount:   200,
						Rate:        4.7,
						Author: &domain.BookAuthor{
							Id:   2,
							Name: "Author 2",
						},
					},
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			books, err := repo.GetAll(ctx)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.expected, books)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
