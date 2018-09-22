package db

import (
	"github.com/asdine/storm"
)

var (
	ErrNotFound      error = storm.ErrNotFound
	ErrAlreadyExists error = storm.ErrAlreadyExists
)
