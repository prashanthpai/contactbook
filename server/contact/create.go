package contact

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/prashanthpai/contactbook/api"
	"github.com/prashanthpai/contactbook/db"

	validator "github.com/asaskevich/govalidator"
)

const createMaxPayloadSize = 1024

func (cb *Book) Create(w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	lr := io.LimitReader(r.Body, createMaxPayloadSize)

	var req api.CreateReq
	if err := json.NewDecoder(lr).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := validator.ValidateStruct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	timestamp := time.Now()
	entry := &db.Entry{
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	if err := cb.db.Store(entry); err != nil {
		if err == db.ErrAlreadyExists {
			http.Error(w, "Contact already exists.", http.StatusConflict)
		} else {
			http.Error(w,
				fmt.Sprintf("Failed to store contact in DB: %s", err.Error()),
				http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", r.URL, req.Email))
	w.WriteHeader(http.StatusCreated)
}
