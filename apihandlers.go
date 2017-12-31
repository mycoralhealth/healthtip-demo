package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"encoding/json"
	"net/http"
)

var auth AuthToken

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
	if err := createAuthToken(user.ID, dbCon); err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, r, http.StatusCreated, auth)
}

func createAuthToken(ID int, dbCon *sql.DB) error {
	// create auth token and put in DB
	auth.Api_user = ID
	// generated salted hash from current time and random number is the auth token
	token, err := getRandomString(10)
	if err != nil {
		return err
	}

	auth.Api_key = hashPassword(token)
	if err := writeAuthToken(dbCon, auth); err != nil {
		return err
	}

	return nil

}

func getRandomString(length int) (string, error) {
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
	r.ParseForm()
	var user User
	user.Email = r.FormValue("email")
	user.Password = hashPassword(r.FormValue("password"))
	if err := checkUserExists(dbCon, user); err != nil {
		handleError(w, r, 404, err.Error())
		return
	}
	ID, err := checkLogin(dbCon, user)
	if err != nil {
		handleError(w, r, 401, err.Error())
		return
	}
	if err := createAuthToken(ID, dbCon); err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, r, http.StatusCreated, auth)

}

func handleLogout(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {

	if err := checkAuthToken(dbCon, auth); err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if err := deleteAuthToken(dbCon, auth); err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	auth = AuthToken{}
	respondWithJSON(w, r, 201, auth)
}

func handleRecords(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
}

func handleSingleRecord(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
}
