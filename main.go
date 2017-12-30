package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func main() {

	dbPath := os.Getenv("DB")
	fmt.Println("opening database: ", dbPath)

	dbCon, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dbCon.Close()

	log.Fatal(run(dbCon))

}

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
	/*
		wrap := func(f func(w http.ResponseWriter, r *http.Request, dbCon *sql.DB)) func(w http.ResponseWriter, r *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {

				f(w, r, dbCon)
			}
		}
	*/
	muxRouter := mux.NewRouter()

	return muxRouter
}
