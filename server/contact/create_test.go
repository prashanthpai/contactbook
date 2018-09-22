package contact

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prashanthpai/contactbook/api"
	"github.com/prashanthpai/contactbook/db"
	"github.com/prashanthpai/contactbook/db/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateNoBody(t *testing.T) {
	assert := require.New(t)

	req, err := http.NewRequest("POST", "/contacts", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	http.HandlerFunc(cb.Create).ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusBadRequest)
	mockDB.AssertNotCalled(t, "Store")
}

func TestCreateLargePayload(t *testing.T) {
	assert := require.New(t)

	createReq := &api.CreateReq{
		Name:  randString(550),
		Email: randString(550),
		Phone: "987654321",
	}

	b, err := json.Marshal(createReq)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/contacts", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	http.HandlerFunc(cb.Create).ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusBadRequest)
	mockDB.AssertNotCalled(t, "Store")
}

func TestCreateInvalidJSON(t *testing.T) {
	assert := require.New(t)

	b := []byte("invlalid json")

	req, err := http.NewRequest("POST", "/contacts", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	http.HandlerFunc(cb.Create).ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusBadRequest)
	mockDB.AssertNotCalled(t, "Store")
}

func TestCreateInvalidReq(t *testing.T) {
	assert := require.New(t)

	createReq := &api.CreateReq{
		Name:  "test name",
		Email: "invalid email",
		Phone: "987654321",
	}

	b, err := json.Marshal(createReq)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/contacts", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	http.HandlerFunc(cb.Create).ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusBadRequest)
	mockDB.AssertNotCalled(t, "Store")
}

func TestCreateContactExists(t *testing.T) {

	assert := require.New(t)

	createReq := &api.CreateReq{
		Name:  "test name",
		Email: "test@test.com",
		Phone: "987654321",
	}

	b, err := json.Marshal(createReq)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/contacts", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	mockDB.On("Store", mock.Anything).Return(db.ErrAlreadyExists)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	http.HandlerFunc(cb.Create).ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusConflict)
	mockDB.AssertCalled(t, "Store", mock.Anything)
}

func TestCreateDBError(t *testing.T) {

	assert := require.New(t)

	createReq := &api.CreateReq{
		Name:  "test name",
		Email: "test@test.com",
		Phone: "987654321",
	}

	b, err := json.Marshal(createReq)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/contacts", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	mockDB.On("Store", mock.Anything).Return(errors.New("some db error"))
	cb := New(mockDB)

	w := httptest.NewRecorder()
	http.HandlerFunc(cb.Create).ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusInternalServerError)
	mockDB.AssertCalled(t, "Store", mock.Anything)
}

func TestCreateSuccess(t *testing.T) {

	assert := require.New(t)

	createReq := &api.CreateReq{
		Name:  "test name",
		Email: "test@test.com",
		Phone: "987654321",
	}

	b, err := json.Marshal(createReq)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/contacts", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	mockDB := new(mocks.DB)
	mockDB.On("Store", mock.Anything).Return(nil)
	cb := New(mockDB)

	w := httptest.NewRecorder()
	http.HandlerFunc(cb.Create).ServeHTTP(w, req)

	assert.Equal(w.Code, http.StatusCreated)
	assert.Equal(w.Result().Header.Get("Location"),
		fmt.Sprintf("/contacts/%s", createReq.Email))
	mockDB.AssertCalled(t, "Store", mock.Anything)
}
