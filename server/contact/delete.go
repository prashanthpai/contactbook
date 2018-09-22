package contact

import (
	"fmt"
	"net/http"

	"github.com/prashanthpai/contactbook/db"

	validator "github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

func (cb *Book) Delete(w http.ResponseWriter, r *http.Request) {

	email := mux.Vars(r)["email"]
	if !validator.IsEmail(email) {
		http.Error(w, "Invalid email.", http.StatusBadRequest)
		return
	}

	if err := cb.db.Delete(email); err != nil {
		if err == db.ErrNotFound {
			http.Error(w, "Contact not found.",
				http.StatusNotFound)
		} else {
			http.Error(w,
				fmt.Sprintf("Failed to delete contact from DB: %s", err.Error()),
				http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
