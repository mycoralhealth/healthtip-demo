package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
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

	// create auth token and put in DB
	var auth AuthToken
	auth.Api_user = user.ID
	// generated salted hash from current time and random number is the auth token
	token, err := getToken(10)
	if err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
	}

	auth.Api_key = hashPassword(token)
	if err := writeAuthToken(dbCon, auth); err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, r, http.StatusCreated, auth)
}

func getToken(length int) (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length], nil
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
