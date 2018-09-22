package contact

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prashanthpai/contactbook/db"
	"github.com/prashanthpai/contactbook/db/mocks"
	"github.com/stretchr/testify/require"
)

func TestReadQueryParams(t *testing.T) {

	assert := require.New(t)

	// multiple query params
	req1, err := http.NewRequest("GET", "/contacts?email=abc@gmail.com&email=xyz@gmail.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	// invalid page number
	req2, err := http.NewRequest("GET", "/contacts?page=-1", nil)
	if err != nil {
		t.Fatal(err)
	}

	req3, err := http.NewRequest("GET", "/contacts?page=notint", nil)
	if err != nil {
		t.Fatal(err)
	}

	reqs := []*http.Request{req1, req2, req3}
	for _, req := range reqs {
		mockDB := new(mocks.DB)
		cb := New(mockDB)

		w := httptest.NewRecorder()
		http.HandlerFunc(cb.Read).ServeHTTP(w, req)

		assert.Equal(w.Code, http.StatusBadRequest)
		mockDB.AssertNotCalled(t, "FindByEmail")
		mockDB.AssertNotCalled(t, "FindByName")
		mockDB.AssertNotCalled(t, "All")
	}
}

func TestReadContactNotFound(t *testing.T) {
	assert := require.New(t)

	email := "xyz@gmail.com"
	req, err := http.NewRequest("GET", fmt.Sprintf("/contacts?email=%s", email), nil)
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	mockDB.On("FindByEmail", email).Return(nil, db.ErrNotFound)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	http.HandlerFunc(cb.Read).ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusNotFound)
	mockDB.AssertCalled(t, "FindByEmail", email)
	mockDB.AssertNotCalled(t, "FindByName")
	mockDB.AssertNotCalled(t, "All")
}

func TestReadDBError(t *testing.T) {
	assert := require.New(t)

	email := "xyz@gmail.com"
	req, err := http.NewRequest("GET", fmt.Sprintf("/contacts?email=%s", email), nil)
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	mockDB.On("FindByEmail", email).Return(nil, errors.New("some db error"))
	cb := New(mockDB)

	w := httptest.NewRecorder()
	http.HandlerFunc(cb.Read).ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusInternalServerError)
	mockDB.AssertCalled(t, "FindByEmail", email)
	mockDB.AssertNotCalled(t, "FindByName")
	mockDB.AssertNotCalled(t, "All")
}

func TestFindByEmailSuccess(t *testing.T) {
	assert := require.New(t)

	email := "xyz@gmail.com"
	req, err := http.NewRequest("GET", fmt.Sprintf("/contacts?email=%s", email), nil)
	if err != nil {
		t.Fatal(err)
	}

	entry := &db.Entry{
		Name:  "name",
		Email: email,
	}

	mockDB := new(mocks.DB)
	mockDB.On("FindByEmail", email).Return(entry, nil)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	http.HandlerFunc(cb.Read).ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusOK)
	mockDB.AssertCalled(t, "FindByEmail", email)
	mockDB.AssertNotCalled(t, "FindByName")
	mockDB.AssertNotCalled(t, "All")
}

func TestFindByNameSuccess(t *testing.T) {
	assert := require.New(t)

	name := "name"
	req, err := http.NewRequest("GET", fmt.Sprintf("/contacts?name=%s&page=1", name), nil)
	if err != nil {
		t.Fatal(err)
	}

	entries := []*db.Entry{
		&db.Entry{Name: "name", Email: "email1"},
		&db.Entry{Name: "name", Email: "email2"},
	}

	mockDB := new(mocks.DB)
	mockDB.On("FindByName", name, 1).Return(entries, nil)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	http.HandlerFunc(cb.Read).ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusOK)
	mockDB.AssertCalled(t, "FindByName", name, 1)
	mockDB.AssertNotCalled(t, "FindByEmail")
	mockDB.AssertNotCalled(t, "All")
}
