package repository

import (
	models "book-api/pkg/model"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func Test_CreateBook(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewBookPostgres(db)

	type mockBehavior func(args models.BookRequest)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         models.BookRequest
		expectedErr  bool
	}{
		{
			name: "OK",
			args: models.BookRequest{
				Name:   "Test book",
				Price:  21,
				Amount: 1,
				Genre:  1,
			},
			expectedErr: false,
			mockBehavior: func(args models.BookRequest) {

				// rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				query := fmt.Sprintf("INSERT INTO %s (.+) RETURNING id", bookTable)
				mock.ExpectQuery(query).
					WithArgs(args.Name, args.Price, args.Genre, args.Amount).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

			},
		},
		{
			name: "Incorrect values",
			args: models.BookRequest{
				Name:  "Test book",
				Price: -21,
				Genre: 1,
			},
			expectedErr: true,
			mockBehavior: func(args models.BookRequest) {

				query := fmt.Sprintf("INSERT INTO %s (.+) RETURNING id", bookTable)
				mock.ExpectQuery(query).
					WithArgs(args.Name, args.Price, args.Genre).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).RowError(1, errors.New("Except argument")))

			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.CreateBook(testCase.args)
			log.Println(got)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, got)
			}

		})
	}
}

func Test_UpdateBook(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewBookPostgres(db)

	type mockBehavior func(args models.Book)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         models.Book
		expectedErr  bool
	}{
		{
			name: "OK",
			args: models.Book{
				Id: 12,
				BookRequest: models.BookRequest{
					Name:   "Test book",
					Price:  100,
					Amount: 23,
					Genre:  1,
				},
			},
			expectedErr: false,
			mockBehavior: func(args models.Book) {
				query := fmt.Sprintf("UPDATE %s (.+) RETURNING id", bookTable)
				mock.ExpectQuery(query).
					WithArgs(args.Name, args.Price, args.Genre, args.Amount, args.Id).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

			},
		},
		{
			name: "Incorrect genre",
			args: models.Book{
				Id: 12,
				BookRequest: models.BookRequest{
					Name:   "Test book",
					Price:  100,
					Amount: 23,
					Genre:  -1,
				},
			},
			expectedErr: true,
			mockBehavior: func(args models.Book) {

				query := fmt.Sprintf("UPDATE %s (.+) RETURNING id", bookTable)
				mock.ExpectQuery(query).
					WithArgs(args.Name, args.Price, args.Genre, args.Amount, args.Id).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).RowError(1, errors.New("Can't update to unkonown genre")))

			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.UpdateBook(testCase.args)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, got)
			}

		})
	}
}

func Test_DeleteBook(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewBookPostgres(db)

	type args struct {
		id int
	}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expectedErr  bool
	}{
		{
			name: "OK",
			args: args{
				id: 1,
			},
			expectedErr: false,
			mockBehavior: func(args args) {
				query := fmt.Sprintf("DELETE FROM %s (.+)", bookTable)
				mock.ExpectQuery(query).
					WithArgs(args.id).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

			},
		},
		{
			name: "Delete by incorrcet id",
			args: args{
				id: -100,
			},
			expectedErr: true,
			mockBehavior: func(args args) {
				query := fmt.Sprintf("DELETE FROM %s (.+)", bookTable)
				mock.ExpectQuery(query).
					WithArgs(args.id).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).RowError(1, errors.New("Can't update to unkonown genre")))

			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			err := r.DeleteBook(testCase.args.id)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func Test_GetBookById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewBookPostgres(db)

	type args struct {
		id int
	}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expectedErr  bool
	}{
		{
			name: "OK",
			args: args{
				id: 1,
			},
			expectedErr: false,
			mockBehavior: func(args args) {
				query := fmt.Sprintf("SELECT id, name, genre, price, amount FROM %s ", bookTable)
				mock.ExpectQuery(query).
					WithArgs(args.id).
					WillReturnRows(sqlmock.NewRows([]string{"id","name","genre","price","amount"}).AddRow(1,"test",1,2,3))

			},
		},
		{
			name: "Empty return",
			args: args{
				id: 1,
			},
			expectedErr: true,
			mockBehavior: func(args args) {
				
				query := fmt.Sprintf("SELECT id, name, genre, price, amount FROM %s ", bookTable)
				mock.ExpectQuery(query).
					WithArgs(args.id).
					WillReturnRows(sqlmock.NewRows([]string{"id","name","genre","price","amount"}))

			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.GetBook(testCase.args.id)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &models.Book{ Id: 1, 
					BookRequest: models.BookRequest{
						Name: "test",
						Genre: 1,
						Price: 2,
						Amount: 3,
					},
				}, got)
			}

		})
	}
}

func Test_GetBooksWithNameFilter(t *testing.T) {
	db, mock, err := sqlmock.Newx(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewBookPostgres(db)

	type args struct {
		bookname string
		}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		isEmpty  bool
	}{
		{
			name: "OK",
			args: args{
				bookname: "test",
			},
			isEmpty: false,
			mockBehavior: func(args args) {
				query := fmt.Sprintf("SELECT id, name, genre, price, amount FROM %s WHERE LOWER(name) LIKE LOWER($1)", bookTable)
				mock.ExpectQuery(query).
					WithArgs(ConvertToLike(args.bookname)).
					WillReturnRows(sqlmock.NewRows([]string{"id","name","genre","price","amount"}).AddRow(1,"test",1,2,3))

			},
		},
		{
			name: "Empty return",
			args: args{
				bookname: "unknkown name",
			},
			isEmpty: true,
			mockBehavior: func(args args) {
				query := fmt.Sprintf("SELECT id, name, genre, price, amount FROM %s WHERE LOWER(name) LIKE LOWER($1)", bookTable)
				mock.ExpectQuery(query).
					WithArgs(ConvertToLike(args.bookname)).
					WillReturnRows(sqlmock.NewRows([]string{"id","name","genre","price","amount"}))

			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.GetBooksWithNameFilter(testCase.args.bookname)
			if testCase.isEmpty {
				var want []models.Book
				assert.NoError(t, err)				
				assert.Equal(t, &want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &[]models.Book{	
					{ Id: 1, 
						BookRequest: models.BookRequest{
							Name: "test",
							Genre: 1,
							Price: 2,
							Amount: 3,
						},
					},
				}, got)
			}

		})
	}
}


func Test_GetBooksWithGenreFilter(t *testing.T) {
	db, mock, err := sqlmock.Newx(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewBookPostgres(db)

	type args struct {
		genre string
		}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		isEmpty  bool
	}{
		{
			name: "OK",
			args: args{
				genre: "test",
			},
			isEmpty: false,
			mockBehavior: func(args args) {
				query := fmt.Sprintf("SELECT B.id, B.name, B.genre, B.price, B.amount FROM %s AS B INNER JOIN %s AS G ON G.id = B.genre WHERE LOWER(G.name) LIKE LOWER($1)", bookTable, genreTable)
				mock.ExpectQuery(query).
					WithArgs(ConvertToLike(args.genre)).
					WillReturnRows(sqlmock.NewRows([]string{"id","name","genre","price","amount"}).AddRow(1,"test",1,2,3))

			},
		},
		{
			name: "Empty return",
			args: args{
				genre: "teteee",
			},
			isEmpty: true,
			mockBehavior: func(args args) {
				query := fmt.Sprintf("SELECT B.id, B.name, B.genre, B.price, B.amount FROM %s AS B INNER JOIN %s AS G ON G.id = B.genre WHERE LOWER(G.name) LIKE LOWER($1)", bookTable, genreTable)
				mock.ExpectQuery(query).
					WithArgs(ConvertToLike(args.genre)).
					WillReturnRows(sqlmock.NewRows([]string{"id","name","genre","price","amount"}))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.GetBooksWithGenreFilter(testCase.args.genre)
			if testCase.isEmpty {
				var want []models.Book
				assert.NoError(t, err)				
				assert.Equal(t, &want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &[]models.Book{	
					{ Id: 1, 
						BookRequest: models.BookRequest{
							Name: "test",
							Genre: 1,
							Price: 2,
							Amount: 3,
						},
					},
				}, got)
			}

		})
	}
}


func Test_GetBooksWithGeneralFilter(t *testing.T) {
	db, mock, err := sqlmock.Newx(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewBookPostgres(db)

	type args struct {
		genre string
		bookname string
		}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		isEmpty  bool
	}{
		{
			name: "OK",
			args: args{
				genre: "test",
				bookname: "test",
			},
			isEmpty: false,
			mockBehavior: func(args args) {
				query := fmt.Sprintf("SELECT B.id, B.name, B.genre, B.price, B.amount FROM %s AS B INNER JOIN %s AS G ON G.id = B.genre  WHERE LOWER(B.name) LIKE LOWER($1) AND LOWER(G.name) LIKE LOWER($2)", bookTable, genreTable)
				mock.ExpectQuery(query).
					WithArgs(ConvertToLike(args.bookname),ConvertToLike(args.genre)).
					WillReturnRows(sqlmock.NewRows([]string{"id","name","genre","price","amount"}).AddRow(1,"test",1,2,3))

			},
		},
		{
			name: "Empty return",
			args: args{
				genre: "teteee",
				bookname: "incorrect data",
			},
			isEmpty: true,
			mockBehavior: func(args args) {
				query := fmt.Sprintf("SELECT B.id, B.name, B.genre, B.price, B.amount FROM %s AS B INNER JOIN %s AS G ON G.id = B.genre  WHERE LOWER(B.name) LIKE LOWER($1) AND LOWER(G.name) LIKE LOWER($2)", bookTable, genreTable)
				mock.ExpectQuery(query).
					WithArgs(ConvertToLike(args.bookname),ConvertToLike(args.genre)).
					WillReturnRows(sqlmock.NewRows([]string{"id","name","genre","price","amount"}))
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.GetBooksWithGeneralFilter(testCase.args.bookname,testCase.args.genre)
			if testCase.isEmpty {
				var want []models.Book
				assert.NoError(t, err)				
				assert.Equal(t, &want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &[]models.Book{	
					{ Id: 1, 
						BookRequest: models.BookRequest{
							Name: "test",
							Genre: 1,
							Price: 2,
							Amount: 3,
						},
					},
				}, got)
			}

		})
	}
}


func Test_GetAllBooks(t *testing.T) {
	db, mock, err := sqlmock.Newx(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewBookPostgres(db)

	type mockBehavior func()

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
	}{
		{
			name: "OK",

			mockBehavior: func() {
				query := fmt.Sprintf("SELECT id, name, genre, price, amount FROM %s", bookTable)
				mock.ExpectQuery(query).
					WillReturnRows(sqlmock.NewRows([]string{"id","name","genre","price","amount"}).AddRow(1,"test",1,2,3))

			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			got, err := r.GetAllBook()
	
				assert.NoError(t, err)
				assert.Equal(t, &[]models.Book{	
					{ Id: 1, 
						BookRequest: models.BookRequest{
							Name: "test",
							Genre: 1,
							Price: 2,
							Amount: 3,
						},
					},
				}, got)
	

		})
	}
}
