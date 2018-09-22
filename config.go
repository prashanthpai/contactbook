package main

import "flag"

const (
	defaultRestAddr = ":8080"
	defaultUser     = "ppai"
	defaultPassword = "livpo"
	defaultDBFile   = "contactbook.db"
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
	flag.StringVar(&cfg.User, "user", defaultUser, "Username of HTTP user.")
	flag.StringVar(&cfg.Password, "password", defaultPassword, "Username of HTTP user.")
}
