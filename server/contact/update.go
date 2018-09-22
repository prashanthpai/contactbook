package contact

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prashanthpai/contactbook/api"
	"github.com/prashanthpai/contactbook/db"

	validator "github.com/asaskevich/govalidator"
)

const updateMaxPayloadSize = 1024

func (cb *Book) Update(w http.ResponseWriter, r *http.Request) {

	email := mux.Vars(r)["email"]
	if !validator.IsEmail(email) {
		http.Error(w, "Invalid email.", http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	lr := io.LimitReader(r.Body, createMaxPayloadSize)

	var req api.UpdateReq
	if err := json.NewDecoder(lr).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := validator.ValidateStruct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" && req.Phone == "" && req.Email == "" {
		http.Error(w, "Name or phone or email must be specified.", http.StatusBadRequest)
		return
	}

	entry := &db.Entry{
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		UpdatedAt: time.Now(),
	}

	if err := cb.db.Update(email, entry); err != nil {
		if err == db.ErrNotFound {
			http.Error(w, "Contact not found.",
				http.StatusNotFound)
		} else {
			http.Error(w,
				fmt.Sprintf("Failed to store contact in DB: %s", err.Error()),
				http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
