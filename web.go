package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
	"github.com/gorilla/mux"
)

func run(dbCon *sql.DB) error {

	httpAddr := os.Getenv("ADDR")

	mux := makeMuxRouter(dbCon)

	log.Println("Listening on ", httpAddr)

	c := cors.New(cors.Options{
		Debug: true,
		AllowedHeaders: []string{"*"},
		AllowedOrigins: []string{"*"}, // All origins. TODO: we need to make this an env and put localhost in here
		AllowedMethods: []string{"GET", "POST", "PUT"}, // Allowing GET, POST, PUT
  	})	

	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        c.Handler(mux),
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

	apiAuth := func(f func(w http.ResponseWriter, r *http.Request, dbCon *sql.DB)) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			apiToken, err := getBasicAPIAuth(r);
			if err != nil {
				handleError(w, r, http.StatusUnauthorized, err.Error())
				return
			}

			if err := checkAPIAuth(dbCon, apiToken); err != nil {
				handleError(w, r, http.StatusUnauthorized, "Unauthorized")
				return
			}

			f(w, r, dbCon)
		}
	}

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/users", wrap(handleWriteUser)).Methods("POST")
	muxRouter.HandleFunc("/users", wrap(handleUpdateUser)).Methods("PUT")
	muxRouter.HandleFunc("/login", wrap(handleLogin)).Methods("POST")
	muxRouter.HandleFunc("/api/logout", apiAuth(handleLogout)).Methods("POST")
	muxRouter.HandleFunc("/api/records", apiAuth(handleRecords)).Methods("GET")
	muxRouter.HandleFunc("/api/records", apiAuth(handleRecords)).Methods("POST")
	muxRouter.HandleFunc("/api/records/{id:[0-9]+}", apiAuth(handleSingleRecord)).Methods("GET")
	muxRouter.HandleFunc("/api/records/{id:[0-9]+}", apiAuth(handleSingleRecord)).Methods("PUT")
	muxRouter.HandleFunc("/api/records/{id:[0-9]+}", apiAuth(handleSingleRecord)).Methods("DELETE")

	return muxRouter
}

func handleError(w http.ResponseWriter, r *http.Request, code int, long string) {
	http.Error(w, http.StatusText(code), code)
	log.Println(r.Method, r.URL, ": HTTP", code, ": ", long)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	//w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
