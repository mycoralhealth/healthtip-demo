package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func run(dbCon *sql.DB) error {

	httpAddr := os.Getenv("ADDR")

	mux := makeMuxRouter(dbCon)

	log.Println("Listening on ", httpAddr)
	s := &http.Server{
		Addr:           httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter(dbCon *sql.DB) http.Handler {
	wrap := func(f func(w http.ResponseWriter, r *http.Request, dbCon *sql.DB)) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {

			f(w, r, dbCon)
		}
	}

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/users", wrap(handleWriteUser)).Methods("POST")
	muxRouter.HandleFunc("/users", wrap(handleUpdateUser)).Methods("PUT")
	muxRouter.HandleFunc("/login", wrap(handleLogin)).Methods("GET")
	muxRouter.HandleFunc("/api/logout", wrap(handleLogout)).Methods("POST")
	muxRouter.HandleFunc("/api/records", wrap(handleRecords)).Methods("GET")
	muxRouter.HandleFunc("/api/records", wrap(handleRecords)).Methods("POST")
	muxRouter.HandleFunc("/api/records/{id:[0-9]+}", wrap(handleSingleRecord)).Methods("GET")
	muxRouter.HandleFunc("/api/records/{id:[0-9]+}", wrap(handleSingleRecord)).Methods("PUT")
	muxRouter.HandleFunc("/api/records/{id:[0-9]+}", wrap(handleSingleRecord)).Methods("DELETE")

	return muxRouter
}
