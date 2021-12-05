package repository

import (
	models "book-api/pkg/model"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type BookPostgres struct {
	db *sqlx.DB
}

func NewBookPostgres(db *sqlx.DB) *BookPostgres {
	return &BookPostgres{ db: db }
}

func (r *BookPostgres) CreateBook(book models.BookRequest) (int, error){
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, price, genre, amount) VALUES ($1,$2,$3,$4) RETURNING id", bookTable)
	
	row := r.db.QueryRow(query, book.Name, book.Price, book.Genre, book.Amount)
	
	if err := row.Err(); err != nil{
		return 0, err 
	}

	if err:= row.Scan(&id); err != nil {
		return 0, err 
	}

	return id, nil
}

func (r *BookPostgres) UpdateBook(book models.Book) (int, error){
	var id int
	query := fmt.Sprintf("UPDATE %s SET name=$1, price=$2, genre=$3, amount=$4 WHERE id = $5 RETURNING id", bookTable)
	
	row := r.db.QueryRow(query, book.Name, book.Price, book.Genre, book.Amount, book.Id)
	
	if err := row.Err(); err != nil{
		return 0, err 
	}

	if err:= row.Scan(&id); err != nil {
		return 0, err 
	}

	return id, nil
}

func (r *BookPostgres) GetBook(id int) (*models.Book, error){
	var rsp models.Book;
	query := fmt.Sprintf("SELECT id, name, genre, price, amount FROM %s WHERE id = $1 LIMIT 1", bookTable)
	err := r.db.Get(&rsp, query, id)
	log.Println(query)
	if err != nil{
		return &rsp, err 
	}

	return &rsp, nil
}
func ConvertToLike(n string) string{
	return "%" + n + "%"
}

func (r *BookPostgres) GetAllBook() (*[]models.Book, error){
	var rsp []models.Book;
	query := fmt.Sprintf("SELECT id, name, genre, price, amount FROM %s", bookTable)
	err := r.db.Select(&rsp, query)
	log.Println(query)
	if err != nil{
		return &rsp, err 
	}
	
	return &rsp, nil
}

func (r *BookPostgres) GetBooksWithNameFilter(name string) (*[]models.Book, error){
	var rsp []models.Book;
	query := fmt.Sprintf("SELECT id, name, genre, price, amount FROM %s WHERE LOWER(name) LIKE LOWER($1)", bookTable)
	err := r.db.Select(&rsp, query, ConvertToLike(name))
	log.Println(query)
	if err != nil{
		return &rsp, err 
	}
	
	return &rsp, nil
}

func (r *BookPostgres) GetBooksWithGenreFilter(name string) (*[]models.Book, error){
	var rsp []models.Book;
	query := fmt.Sprintf("SELECT B.id, B.name, B.genre, B.price, B.amount FROM %s AS B INNER JOIN %s AS G ON G.id = B.genre WHERE LOWER(G.name) LIKE LOWER($1)", bookTable, genreTable)
	err := r.db.Select(&rsp, query, ConvertToLike(name))
	log.Println(query)
	if err != nil{
		return &rsp, err 
	}
	
	return &rsp, nil
}

func (r *BookPostgres) GetBooksWithGeneralFilter(bookname string, genre string) (*[]models.Book, error){
	var rsp []models.Book;
	query := fmt.Sprintf("SELECT B.id, B.name, B.genre, B.price, B.amount FROM %s AS B INNER JOIN %s AS G ON G.id = B.genre  WHERE LOWER(B.name) LIKE LOWER($1) AND LOWER(G.name) LIKE LOWER($2)", bookTable, genreTable)
	err := r.db.Select(&rsp, query, ConvertToLike(bookname), ConvertToLike(genre))
	log.Println(query)
	if err != nil{
		return &rsp, err 
	}
	
	return &rsp, nil
}

func (r *BookPostgres) DeleteBook(id int) error{
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", bookTable)
	row := r.db.QueryRow(query, id)
	
	if err := row.Err(); err != nil{
		return err 
	}

	if err:= row.Scan(&id); err != nil {
		return err 
	}

	return nil
}