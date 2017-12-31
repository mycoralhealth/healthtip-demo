package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func handleWriteUser(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	var user User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		handleError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	lastID, err := writeUser(dbCon, user)
	if err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	user.ID = int(lastID)
	respondWithJSON(w, r, http.StatusCreated, user)
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	var user User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		handleError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := updateUser(dbCon, user); err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, r, http.StatusCreated, user)

}

func handleLogin(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
}
func handleLogout(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
}
func handleRecords(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
}
func handleSingleRecord(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
}
