package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prashanthpai/contactbook/db"
	"github.com/prashanthpai/contactbook/server"
	"github.com/prashanthpai/contactbook/server/auth"
	"github.com/prashanthpai/contactbook/server/contact"
)

func main() {
	flag.Parse()

	// initialize and open DB
	db, err := db.New(cfg.DBFile)
	if err != nil {
		log.Fatalf("db.New(%s) failed: %s", cfg.DBFile, err)
	}
	defer db.Close()

	log.Printf("Opened DB file %s", cfg.DBFile)

	cb := contact.New(db)

	authconfig := &auth.Config{
		User:     cfg.User,
		Password: cfg.Password,
	}

	// start HTTP server
	srv := server.New(cfg.Addr, cb, authconfig)
	go func(s *http.Server) {
		log.Printf("Starting HTTP service on address: %s", cfg.Addr)
		if err := s.ListenAndServe(); err != nil {
			log.Fatalf("http.Server.ListenAndServe() failed: %s\n", err)
		}
	}(srv)

	// block main goroutine for interrupt
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
	log.Println("Received interrupt. Shutting down...")

	// gracefully shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("http.Server.Shutdown() failed: %s\n", err)
	}
}
