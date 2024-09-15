package domain

type BookEntity struct {
	Id int
	BookInfo
	Genres []BookGenre
}

type BookInfo struct {
	Author      *BookAuthor
	Cover       []byte
	Description string
	Name        string
	PageCount   int
	Rate        float32
}

type BookAuthor struct {
	Id    int
	Name  string
	Books []int
}

type GenreBookEntity struct {
	GenreId int
	BookId  int
}
type BookGenre struct {
	Id   int
	Name string
}
