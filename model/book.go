package models

type Book struct {
	Id int `json:"id"`
	BookRequest
}

type BookRequest struct {
	Name   string  `json:"name" validate:"required,max=100"`
	Price  float32 `json:"price" validate:"required,min=0"`
	Genre  int     `json:"genre" validate:"required"`
	Amount int     `json:"amount" validate:"required,min=0"`
}

type SimpleResponse struct {
	Id int `json:"id"`
}
