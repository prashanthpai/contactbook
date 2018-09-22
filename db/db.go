package db

import (
	"time"
)

type Entry struct {
	ID        int    `storm:"id,increment"`
	Email     string `storm:"unique"`
	Name      string `storm:"index"`
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DB interface {
	Store(*Entry) error
	Update(string, *Entry) error
	All(int) ([]*Entry, error)
	FindByEmail(string) (*Entry, error)
	FindByName(string, int) ([]*Entry, error)
	Delete(string) error
	Close() error
}
