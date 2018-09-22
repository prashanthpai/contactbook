// +build integration

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/prashanthpai/contactbook/api"
	"github.com/stretchr/testify/require"
)

func TestContactBook(t *testing.T) {

	p, err := spawnProcess()
	if err != nil {
		t.Fatal(err)
	}
	defer p.Destroy()

	//TODO: Several negative tests haven't been tested yet.

	t.Run("testCreate", func(t *testing.T) { testCreate(t) })
	t.Run("testRead", func(t *testing.T) { testRead(t) })
	t.Run("testUpdate", func(t *testing.T) { testUpdate(t) })
	t.Run("testDelete", func(t *testing.T) { testDelete(t) })
}

func testCreate(t *testing.T) {

	assert := require.New(t)

	client := newClient()

	for i := 0; i < 25; i++ {
		createReq := api.CreateReq{
			Name:  "Bruce Wayne",
			Email: fmt.Sprintf("bwayne%d@example.com", i+1),
			Phone: randNum(10),
		}

		b, err := json.Marshal(createReq)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "http://localhost:8080/contacts", bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("content-type", "application/json")
		req.SetBasicAuth(testUser, testPassword)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(resp.StatusCode, http.StatusCreated)
	}

	// test existing contact
	createReq := api.CreateReq{
		Name:  "New Name",
		Email: "bwayne1@example.com",
		Phone: randNum(10),
	}

	b, err := json.Marshal(createReq)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/contacts", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth(testUser, testPassword)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(resp.StatusCode, http.StatusConflict)
}

func testRead(t *testing.T) {

	assert := require.New(t)

	client := newClient()

	var contacts api.Contacts
	var url string

	for i := 0; i <= 4; i++ {
		if i == 0 {
			url = "http://localhost:8080/contacts"
		} else {
			url = fmt.Sprintf("http://localhost:8080/contacts?page=%d", i)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetBasicAuth(testUser, testPassword)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if i <= 3 {
			assert.Equal(resp.StatusCode, http.StatusOK)
		} else {
			assert.Equal(resp.StatusCode, http.StatusNotFound)
			continue
		}

		if err := json.NewDecoder(resp.Body).Decode(&contacts); err != nil {
			t.Fatal(err)
		}

		if i < 3 {
			assert.Equal(len(contacts), 10)
		} else if i == 3 {
			assert.Equal(len(contacts), 5)
		}
	}
}

func testUpdate(t *testing.T) {

	assert := require.New(t)

	client := newClient()

	for i := 0; i < 10; i++ {
		updateReq := api.UpdateReq{
			Name:  "Batman",
			Phone: randNum(10),
		}

		b, err := json.Marshal(updateReq)
		if err != nil {
			t.Fatal(err)
		}

		email := fmt.Sprintf("bwayne%d@example.com", i+1)
		req, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:8080/contacts/%s", email), bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("content-type", "application/json")
		req.SetBasicAuth(testUser, testPassword)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(resp.StatusCode, http.StatusOK)
	}

	// ensure that update was successful
	for i := 0; i < 10; i++ {
		email := fmt.Sprintf("bwayne%d@example.com", i+1)
		url := fmt.Sprintf("http://localhost:8080/contacts?email=%s", email)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetBasicAuth(testUser, testPassword)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		assert.Equal(resp.StatusCode, http.StatusOK)

		var contacts api.Contacts
		if err := json.NewDecoder(resp.Body).Decode(&contacts); err != nil {
			t.Fatal(err)
		}

		assert.Equal(len(contacts), 1)
		contact := contacts[0]
		assert.Equal(contact.Email, email)
		assert.Equal(contact.Name, "Batman")
	}
}

func testDelete(t *testing.T) {

	assert := require.New(t)

	client := newClient()

	for i := 0; i < 25; i++ {

		email := fmt.Sprintf("bwayne%d@example.com", i+1)
		url := fmt.Sprintf("http://localhost:8080/contacts/%s", email)
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetBasicAuth(testUser, testPassword)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(resp.StatusCode, http.StatusNoContent)
	}

	// ensure that delete was successful
	for i := 0; i < 25; i++ {

		email := fmt.Sprintf("bwayne%d@example.com", i+1)
		url := fmt.Sprintf("http://localhost:8080/contacts?email=%s", email)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetBasicAuth(testUser, testPassword)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		assert.Equal(resp.StatusCode, http.StatusNotFound)
	}
}
