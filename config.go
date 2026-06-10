package main

import (
	"flag"
	"os"
)

const (
	defaultRestAddr = ":8080"
	defaultDBFile   = "contactbook.db"
	defaultUserEnv  = "CONTACTBOOK_USER"
	defaultPassEnv  = "CONTACTBOOK_PASSWORD"
)

type config struct {
	Addr     string
	DBFile   string
	User     string
	Password string
}

var cfg *config

func init() {
	cfg = new(config)

	flag.StringVar(&cfg.Addr, "addr", defaultRestAddr, "Address to listen on for HTTP server.")
	flag.StringVar(&cfg.DBFile, "db-file", defaultDBFile, "Path to db file.")
	flag.StringVar(&cfg.User, "user", os.Getenv(defaultUserEnv), "Username of HTTP user.")
	flag.StringVar(&cfg.Password, "password", os.Getenv(defaultPassEnv), "Password of HTTP user.")
}
