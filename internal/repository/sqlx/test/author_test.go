package sqlx_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"testing"

	"book-api/internal/domain"
	"book-api/internal/repository/sqlx"
	"book-api/internal/transaction"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	provider "github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdateAuthor(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := provider.NewDb(db, "sqlmock")
	txManager := transaction.NewTransactionManager(sqlxDB)
	repo := sqlx.NewSQLXAuthorRepository(sqlxDB)
	ctx := context.Background()

	tests := []struct {
		testName  string
		mockSetup func()

		author     domain.BookAuthor
		expectErr  bool
		expectedID int
	}{
		{
			testName: "Insert new author",
			author:   domain.BookAuthor{Name: "John Doe"},
			mockSetup: func() {
				query := fmt.Sprintf(`INSERT INTO %s \(name\) VALUES \(\$1\) ON CONFLICT \(name\) DO UPDATE SET name = EXCLUDED.name; RETURNING id`, sqlx.AuthorTable)
				mock.ExpectBegin()
				mock.ExpectQuery(query).WithArgs("John Doe").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectErr:  false,
			expectedID: 1,
		},
		{
			testName: "Update existing author",
			author:   domain.BookAuthor{Name: "Jane Doe"},
			mockSetup: func() {
				query := fmt.Sprintf(`INSERT INTO %s \(name\) VALUES \(\$1\) ON CONFLICT \(name\) DO UPDATE SET name = EXCLUDED.name; RETURNING id`, sqlx.AuthorTable)
				mock.ExpectBegin()
				mock.ExpectQuery(query).WithArgs("Jane Doe").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
				mock.ExpectCommit()
			},
			expectErr:  false,
			expectedID: 2,
		},
		{
			testName: "Error during upsert",
			author:   domain.BookAuthor{Name: "Error Author"},
			mockSetup: func() {
				query := fmt.Sprintf(`INSERT INTO %s \(name\) VALUES \(\$1\) ON CONFLICT \(name\) DO UPDATE SET name = EXCLUDED.name; RETURNING id`, sqlx.AuthorTable)
				expectedErr := fmt.Errorf("db error")
				mock.ExpectBegin()
				mock.ExpectQuery(query).WithArgs("Error Author").WillReturnError(expectedErr)
				mock.ExpectRollback().WillReturnError(expectedErr)
			},
			expectErr:  true,
			expectedID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.mockSetup()

			err := txManager.RunInTransaction(ctx, func(tx transaction.Transaction) error {
				return repo.CreateOrUpdate(ctx, tx, &tt.author)
			})

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, tt.author.Id)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetByNameAuthor(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := provider.NewDb(db, "sqlmock")
	repo := sqlx.NewSQLXAuthorRepository(sqlxDB)
	ctx := context.Background()

	tests := []struct {
		testName  string
		mockSetup func()

		authorName     string
		expectedAuthor *domain.BookAuthor
		expectErr      bool
	}{
		{
			testName: "Successfull selection",
			mockSetup: func() {
				query := fmt.Sprintf(`SELECT id, name FROM %s WHERE name = \$1 LIMIT 1`, sqlx.AuthorTable)
				rows := sqlmock.NewRowsWithColumnDefinition(sqlmock.NewColumn("id"), sqlmock.NewColumn("name")).AddRows([]driver.Value{1, "John Doe"})
				mock.ExpectQuery(query).WithArgs("John Doe").WillReturnRows(rows)
			},

			authorName:     "John Doe",
			expectedAuthor: &domain.BookAuthor{Name: "John Doe", Id: 1},
			expectErr:      false,
		},
		{
			testName: "Selected item doesn't exist",
			mockSetup: func() {
				query := fmt.Sprintf(`SELECT id, name FROM %s WHERE name = \$1 LIMIT 1`, sqlx.AuthorTable)
				rows := sqlmock.NewRowsWithColumnDefinition(sqlmock.NewColumn("id"), sqlmock.NewColumn("name"))
				mock.ExpectQuery(query).WithArgs("John Doe").WillReturnRows(rows)
			},

			authorName:     "John Doe",
			expectedAuthor: &domain.BookAuthor{Name: "John Doe", Id: 1},
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.mockSetup()

			actualAuthor, err := repo.GetByName(ctx, tt.authorName)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAuthor, actualAuthor)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
