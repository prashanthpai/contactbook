package contact

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/prashanthpai/contactbook/db"
	"github.com/prashanthpai/contactbook/db/mocks"
	"github.com/stretchr/testify/require"
)

func newDeleteHandler(f http.HandlerFunc) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/contacts/{email}", f).Methods("DELETE")
	return r
}

func TestDeleteInvalidEmail(t *testing.T) {

	assert := require.New(t)

	req, err := http.NewRequest("DELETE", "/contacts/invalidemail", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	handler := newDeleteHandler(cb.Delete)
	handler.ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusBadRequest)
	mockDB.AssertNotCalled(t, "Delete")
}

func TestDeleteNotFound(t *testing.T) {

	assert := require.New(t)

	email := "non_existent_email@gmail.com"
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/contacts/%s", email), nil)
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	mockDB.On("Delete", email).Return(db.ErrNotFound)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	handler := newDeleteHandler(cb.Delete)
	handler.ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusNotFound)
	mockDB.AssertCalled(t, "Delete", email)
}

func TestDeleteDBError(t *testing.T) {

	assert := require.New(t)

	email := "valid_email@gmail.com"
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/contacts/%s", email), nil)
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	mockDB.On("Delete", email).Return(errors.New("some error"))
	cb := New(mockDB)

	w := httptest.NewRecorder()
	handler := newDeleteHandler(cb.Delete)
	handler.ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusInternalServerError)
	mockDB.AssertCalled(t, "Delete", email)
}

func TestDeleteSuccess(t *testing.T) {

	assert := require.New(t)

	email := "valid_email@gmail.com"
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/contacts/%s", email), nil)
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	mockDB.On("Delete", email).Return(nil)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	handler := newDeleteHandler(cb.Delete)
	handler.ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusNoContent)
	mockDB.AssertCalled(t, "Delete", email)
}
