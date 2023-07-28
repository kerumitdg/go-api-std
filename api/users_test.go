package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fredrikaverpil/go-api-std/stores"
)

func TestCreateUserOk(t *testing.T) {
	expectedJsonBody := `{"id":1, "username":"john"}`

	store := stores.NewDummyStore()
	server := NewServer(":8080", &store)
	url := "/users"
	body := bytes.NewBufferString(`{"username":"john", "password":"secret"}`)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	server.router.ServeHTTP(rr, req)

	assert.Exactly(t, rr.Code, http.StatusCreated)
	assert.JSONEq(t, expectedJsonBody, rr.Body.String())
}

func TestCreateUserNoUsername(t *testing.T) {
	store := stores.NewDummyStore()
	server := NewServer(":8080", &store)
	url := "/users"
	body := bytes.NewBufferString(`{"password":"secret"}`)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	server.router.ServeHTTP(rr, req)

	assert.Exactly(t, rr.Code, http.StatusBadRequest)
}

func TestCreateUserNoPassword(t *testing.T) {
	store := stores.NewDummyStore()
	server := NewServer(":8080", &store)
	url := "/users"
	body := bytes.NewBufferString(`{"username":"john"}`)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	server.router.ServeHTTP(rr, req)

	assert.Exactly(t, rr.Code, http.StatusBadRequest)
}

func TestCreateUserUsernameTaken(t *testing.T) {
	store := stores.NewDummyStore()
	server := NewServer(":8080", &store)
	url := "/users"
	body := bytes.NewBufferString(`{"username":"foo", "password":"secret"}`)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	server.router.ServeHTTP(rr, req)

	assert.Exactly(t, rr.Code, http.StatusCreated)

	// User 2, with the same username
	body = bytes.NewBufferString(`{"username":"foo", "password":"secret"}`)
	req, err = http.NewRequest("POST", url, body)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()

	server.router.ServeHTTP(rr, req)

	assert.Exactly(t, rr.Code, http.StatusConflict)
}

func TestGetUserByIdOk(t *testing.T) {
	expectedJsonBody := `{"id":1, "username":"john"}`

	store := stores.NewDummyStore()
	user, _ := store.CreateUser("john", "secret")
	assert.Exactly(t, user.ID, 1)

	server := NewServer(":8080", &store)
	url := "/users/1"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	server.router.ServeHTTP(rr, req)

	assert.Exactly(t, rr.Code, 200)
	assert.JSONEq(t, expectedJsonBody, rr.Body.String())
}

func TestGetUsersNotSupported(t *testing.T) {
	store := stores.NewDummyStore()
	server := NewServer(":8080", &store)
	url := "/users/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	server.router.ServeHTTP(rr, req)

	assert.Exactly(t, rr.Code, http.StatusNotFound)
}

func TestUsersNoSlash(t *testing.T) {
	store := stores.NewDummyStore()
	server := NewServer(":8080", &store)
	url := "/users"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	server.router.ServeHTTP(rr, req)

	assert.Exactly(t, rr.Code, http.StatusMethodNotAllowed)
}

func TestNonExistingUser(t *testing.T) {
	store := stores.NewDummyStore()
	server := NewServer(":8080", &store)
	url := "/users/0"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	server.router.ServeHTTP(rr, req)

	assert.Exactly(t, rr.Code, http.StatusNotFound)
}
