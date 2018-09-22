package server

import (
	"net/http"
	"os"
	"time"

	"github.com/prashanthpai/contactbook/server/auth"
	"github.com/prashanthpai/contactbook/server/contact"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func New(addr string, cb *contact.Book, authconfig *auth.Config) *http.Server {
	router := mux.NewRouter()

	router.HandleFunc("/contacts", cb.Create).Methods("POST")
	router.HandleFunc("/contacts", cb.Read).Methods("GET")
	router.HandleFunc("/contacts/{email}", cb.Update).Methods("PUT")
	router.HandleFunc("/contacts/{email}", cb.Delete).Methods("DELETE")

	chain := handlers.RecoveryHandler()(
		handlers.LoggingHandler(os.Stdout,
			auth.BasicAuth(authconfig, router)),
	)

	return &http.Server{
		Addr:         addr,
		Handler:      chain,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
