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

func TestUpsertGenreToBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := provider.NewDb(db, "sqlmock")
	txManager := transaction.NewTransactionManager(sqlxDB)
	repo := sqlx.NewSQLXAGenreRepository(sqlxDB)
	ctx := context.Background()

	tests := []struct {
		testName  string
		mockSetup func()
		bookID    int
		genreID   int
		expectErr bool
	}{
		{
			testName: "Upsert genre to book successfully",
			bookID:   1,
			genreID:  2,
			mockSetup: func() {
				query := fmt.Sprintf(`INSERT INTO %s \(book_id, genre_id\) VALUES \(\$1, \$2\) ON CONFLICT \(book_id\) DO UPDATE SET genre_id = EXCLUDED.genre_id`, sqlx.BookGenreTable)
				mock.ExpectBegin()
				mock.ExpectExec(query).WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			testName: "Error during upsert",
			bookID:   1,
			genreID:  2,
			mockSetup: func() {
				query := fmt.Sprintf(`INSERT INTO %s \(book_id, genre_id\) VALUES \(\$1, \$2\) ON CONFLICT \(book_id\) DO UPDATE SET genre_id = EXCLUDED.genre_id`, sqlx.BookGenreTable)
				mock.ExpectBegin()
				mock.ExpectExec(query).WithArgs(1, 2).WillReturnError(fmt.Errorf("db error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.mockSetup()

			err := txManager.RunInTransaction(ctx, func(tx transaction.Transaction) error {
				return repo.UpsertGenreToBook(ctx, tx, tt.bookID, tt.genreID)
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

func TestGetBookGenreById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := provider.NewDb(db, "sqlmock")
	repo := sqlx.NewSQLXAGenreRepository(sqlxDB)
	ctx := context.Background()

	tests := []struct {
		testName       string
		mockSetup      func()
		bookID         int
		expectedGenres []domain.BookGenre
		expectErr      bool
	}{
		{
			testName: "Successful retrieval of genres by book ID",
			bookID:   1,
			mockSetup: func() {
				query := fmt.Sprintf(`SELECT G.id AS Id, G.name AS Name FROM %s AS BG INNER JOIN %s AS G ON BG.genre_id = G.id WHERE BG.book_id = \$1;`, sqlx.BookGenreTable, sqlx.GenreTable)
				rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Fiction").AddRow(2, "Science Fiction")
				mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
			},
			expectedGenres: []domain.BookGenre{
				{Id: 1, Name: "Fiction"},
				{Id: 2, Name: "Science Fiction"},
			},
			expectErr: false,
		},
		{
			testName: "No genres found for book ID",
			bookID:   1,
			mockSetup: func() {
				query := fmt.Sprintf(`SELECT G.id AS Id, G.name AS Name FROM %s AS BG INNER JOIN %s AS G ON BG.genre_id = G.id WHERE BG.book_id = \$1;`, sqlx.BookGenreTable, sqlx.GenreTable)
				rows := sqlmock.NewRows([]string{"id", "name"})
				mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
			},
			expectedGenres: []domain.BookGenre{},
			expectErr:      false,
		},
		{
			testName: "Error during retrieval of genres",
			bookID:   1,
			mockSetup: func() {
				query := fmt.Sprintf(`SELECT G.id AS Id, G.name AS Name FROM %s AS BG INNER JOIN %s AS G ON BG.genre_id = G.id WHERE BG.book_id = \$1;`, sqlx.BookGenreTable, sqlx.GenreTable)
				mock.ExpectQuery(query).WithArgs(1).WillReturnError(fmt.Errorf("db error"))
			},
			expectedGenres: nil,
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.mockSetup()

			genres, err := repo.GetByBookId(ctx, tt.bookID)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.expectedGenres, genres)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteBookGenre(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := provider.NewDb(db, "sqlmock")
	txManager := transaction.NewTransactionManager(sqlxDB)
	repo := sqlx.NewSQLXAGenreRepository(sqlxDB)
	ctx := context.Background()

	tests := []struct {
		testName  string
		mockSetup func()
		bookID    int
		expectErr bool
	}{
		{
			testName: "Delete book genre successfully",
			bookID:   1,
			mockSetup: func() {
				query := fmt.Sprintf(`DELETE FROM %s WHERE book_id = \$1`, sqlx.BookGenreTable)
				mock.ExpectBegin()
				mock.ExpectExec(query).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			testName: "Error during deletion",
			bookID:   1,
			mockSetup: func() {
				query := fmt.Sprintf(`DELETE FROM %s WHERE book_id = \$1`, sqlx.BookGenreTable)
				mock.ExpectBegin()
				mock.ExpectExec(query).WithArgs(1).WillReturnError(fmt.Errorf("db error"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.mockSetup()

			err := txManager.RunInTransaction(ctx, func(tx transaction.Transaction) error {
				return repo.DeleteBookGenre(ctx, tx, tt.bookID)
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
