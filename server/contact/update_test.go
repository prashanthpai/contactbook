package contact

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/prashanthpai/contactbook/api"
	"github.com/prashanthpai/contactbook/db"
	"github.com/prashanthpai/contactbook/db/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newUpdateHandler(f http.HandlerFunc) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/contacts/{email}", f).Methods("PUT")
	return r
}

func TestUpdateNoBody(t *testing.T) {
	assert := require.New(t)

	req, err := http.NewRequest("PUT", "/contacts/email@test.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	handler := newUpdateHandler(cb.Update)
	handler.ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusBadRequest)
	mockDB.AssertNotCalled(t, "Update")
}

func TestUpdateInvalidEmail1(t *testing.T) {
	assert := require.New(t)

	req, err := http.NewRequest("PUT", "/contacts/invalidemail", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	handler := newUpdateHandler(cb.Update)
	handler.ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusBadRequest)
	mockDB.AssertNotCalled(t, "Update")
}

func TestUpdateLargePayload(t *testing.T) {
	assert := require.New(t)

	updateReq := &api.UpdateReq{
		Name:  randString(550),
		Email: randString(550),
		Phone: "987654321",
	}

	b, err := json.Marshal(updateReq)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/contacts/email@test.com", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	handler := newUpdateHandler(cb.Update)
	handler.ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusBadRequest)
	mockDB.AssertNotCalled(t, "Update")
}

func TestUpdateInvalidJSON(t *testing.T) {
	assert := require.New(t)

	b := []byte("invlalid json")

	req, err := http.NewRequest("PUT", "/contacts/email@test.com", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	handler := newUpdateHandler(cb.Update)
	handler.ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusBadRequest)
	mockDB.AssertNotCalled(t, "Update")
}

func TestUpdateInvalidReq(t *testing.T) {
	assert := require.New(t)

	updateReq := &api.UpdateReq{
		Name:  "test name",
		Email: "invalid email",
		Phone: "987654321",
	}

	b, err := json.Marshal(updateReq)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/contacts/email@test.com", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	handler := newUpdateHandler(cb.Update)
	handler.ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusBadRequest)
	mockDB.AssertNotCalled(t, "Update")
}

func TestUpdateContactNotExist(t *testing.T) {

	assert := require.New(t)

	updateReq := &api.UpdateReq{
		Name:  "test name",
		Email: "test@test.com",
		Phone: "987654321",
	}

	b, err := json.Marshal(updateReq)
	if err != nil {
		t.Fatal(err)
	}

	email := "non_existent_email@test.com"
	req, err := http.NewRequest("PUT", fmt.Sprintf("/contacts/%s", email), bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	mockDB.On("Update", email, mock.Anything).Return(db.ErrNotFound)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	handler := newUpdateHandler(cb.Update)
	handler.ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusNotFound)
	mockDB.AssertCalled(t, "Update", email, mock.Anything)
}

func TestUpdateDBError(t *testing.T) {

	assert := require.New(t)

	updateReq := &api.UpdateReq{
		Name:  "test name",
		Email: "test@test.com",
		Phone: "987654321",
	}

	b, err := json.Marshal(updateReq)
	if err != nil {
		t.Fatal(err)
	}

	email := "email@test.com"
	req, err := http.NewRequest("PUT", fmt.Sprintf("/contacts/%s", email), bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	mockDB.On("Update", email, mock.Anything).Return(errors.New("some db error"))
	cb := New(mockDB)

	w := httptest.NewRecorder()
	handler := newUpdateHandler(cb.Update)
	handler.ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusInternalServerError)
	mockDB.AssertCalled(t, "Update", email, mock.Anything)
}

func TestUpdateSuccess(t *testing.T) {

	assert := require.New(t)

	updateReq := &api.UpdateReq{
		Name:  "test name",
		Email: "test@test.com",
		Phone: "987654321",
	}

	b, err := json.Marshal(updateReq)
	if err != nil {
		t.Fatal(err)
	}

	email := "email@test.com"
	req, err := http.NewRequest("PUT", fmt.Sprintf("/contacts/%s", email), bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	mockDB.On("Update", email, mock.Anything).Return(nil)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	handler := newUpdateHandler(cb.Update)
	handler.ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusOK)
	mockDB.AssertCalled(t, "Update", email, mock.Anything)
}
