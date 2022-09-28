package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/irmatov/companies/middleware"
	"github.com/irmatov/companies/postgres"
	"github.com/irmatov/companies/server"
)

const countryLookupTimeout = 5 * time.Second

func main() {
	db, err := sql.Open("pgx", os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}
	srv := server.New(postgres.New(db))
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	httpServer := &http.Server{
		Addr: ":8080",
		Handler: &middleware.Country{
			Next:               srv,
			LookupURLFormat:    "https://ipapi.co/%s/json/",
			AllowedCountryCode: os.Getenv("ALLOWED_COUNTRY_CODE"),
			Client:             &http.Client{Timeout: countryLookupTimeout},
		},
	}
	go httpServer.ListenAndServe()
	<-ch
	httpServer.Shutdown(context.Background())
}
