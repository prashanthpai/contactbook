package contact

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/prashanthpai/contactbook/api"
	"github.com/prashanthpai/contactbook/db"
)

func (cb *Book) Read(w http.ResponseWriter, r *http.Request) {

	name, nameOk := r.URL.Query()["name"]
	email, emailOk := r.URL.Query()["email"]
	page, pageOk := r.URL.Query()["page"]

	if len(name) > 1 || len(email) > 1 || len(page) > 1 {
		http.Error(w, "Service supports only one value per query param.", http.StatusBadRequest)
		return
	}

	var err error
	pageNum := 1
	if pageOk {
		pageNum, err = strconv.Atoi(page[0])
		if err != nil || pageNum <= 0 {
			http.Error(w, "Invalid page number.", http.StatusBadRequest)
			return
		}
	}

	// fetch results from the DB
	var results []*db.Entry
	if emailOk {
		results, err = cb.readEntries(email[0], "", pageNum)
	} else if nameOk {
		results, err = cb.readEntries("", name[0], pageNum)
	} else {
		results, err = cb.readEntries("", "", pageNum)
	}
	if err != nil {
		if err == db.ErrNotFound {
			http.Error(w, "Contact not found.", http.StatusNotFound)
		} else {
			http.Error(w,
				fmt.Sprintf("Failed to fetch contacts from DB: %s", err.Error()),
				http.StatusInternalServerError)
		}
		return
	}

	if len(results) == 0 {
		http.Error(w, "Contacts not found.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var contacts api.Contacts
	for _, entry := range results {
		contacts = append(contacts, api.Contact{
			Name:      entry.Name,
			Email:     entry.Email,
			Phone:     entry.Phone,
			CreatedAt: entry.CreatedAt,
			UpdatedAt: entry.UpdatedAt,
		})
	}

	if err := json.NewEncoder(w).Encode(contacts); err != nil {
		log.Printf("Encoding json failed: %s", err)
	}
}

func (cb *Book) readEntries(email string, name string, pageNum int) ([]*db.Entry, error) {

	var results []*db.Entry

	if email != "" {
		entry, err := cb.db.FindByEmail(email)
		if err != nil {
			return nil, err
		}
		results = append(results, entry)
	} else if name != "" {
		entries, err := cb.db.FindByName(name, pageNum)
		if err != nil {
			return nil, err
		}
		results = append(results, entries...)
	} else {
		entries, err := cb.db.All(pageNum)
		if err != nil {
			return nil, err
		}
		results = append(results, entries...)
	}

	return results, nil
}
