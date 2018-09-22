package contact

import (
	"github.com/prashanthpai/contactbook/db"
)

type Book struct {
	db db.DB
}

func New(db db.DB) *Book {
	return &Book{
		db: db,
	}
}
