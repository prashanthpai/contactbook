// Package api provides HTTP request and response structs.
package api

import (
	"time"
)

type Contact struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
}

type Contacts []Contact
