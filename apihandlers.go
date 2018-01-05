package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"errors"

	"github.com/gorilla/mux"
)

func makeLoginResult(user User, auth AuthToken) LoginResult {
	var result LoginResult

	result.Token = auth
	result.Email = user.Email
	result.First_name = user.First_name
	result.Last_name = user.Last_name

	return result
}

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
	auth, err := createAuthToken(user.ID, dbCon)
	if err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, r, http.StatusCreated, makeLoginResult(user, auth))
}

func createAuthToken(ID int, dbCon *sql.DB) (AuthToken, error) {
	var auth AuthToken
	// create auth token and put in DB
	auth.Api_user = ID
	// generated salted hash from current time and random number is the auth token
	token, err := getRandomString(10)
	if err != nil {
		return auth, err
	}

	auth.Api_key = hashPassword(token)
	if err := writeAuthToken(dbCon, auth); err != nil {
		return auth, err
	}

	return auth, nil
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

	respondWithJSON(w, r, http.StatusCreated, user.Email)

}

func handleLogin(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	
	user, err := getBasicLoginAuth(r);

	if err != nil {
		handleError(w, r, http.StatusUnauthorized, err.Error())
		return
	}

	dbUser, err := checkUserExists(dbCon, user); 

	if err != nil {
		handleError(w, r, 404, err.Error())
		return
	}

	if err := checkLoginAuth(dbCon, dbUser); err != nil {
		handleError(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	auth, err := createAuthToken(dbUser.ID, dbCon)
	if err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, r, http.StatusOK, makeLoginResult(dbUser, auth))
}

func getBasicLoginAuth(r *http.Request) (User, error) {
	var user User

	email, password, ok := r.BasicAuth()
	if !ok {
		return user, errors.New("Invalid login auth")
	}

	user.Email = email
	user.Password = password

	return user, nil
}

/*
API handlers, authorization is handled in web.go wrapper function
*/

func getBasicAPIAuth(r *http.Request) (AuthToken, error) {
	var token AuthToken

	p, s, ok := r.BasicAuth()
	if !ok {
		return token, errors.New("Invalid API Authorization token. Required Basic Api_user:Api_key.")
	}

	apiUser, err := strconv.Atoi(p)
	if err != nil {
		return token, errors.New("Invalid Api_user.")
	}

	token.Api_user = apiUser
	token.Api_key = s

	return token, nil
}

func handleLogout(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	auth, _ := getBasicAPIAuth(r);

	if err := deleteAuthToken(dbCon, auth); err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, r, http.StatusNoContent, nil)
}

func handleRecords(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	auth, _ := getBasicAPIAuth(r);

	if r.Method == "GET" {
		records, err := getAllRecords(auth.Api_user, dbCon)
		if err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		bytes, err := json.MarshalIndent(records, "", "  ")
		if err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		io.WriteString(w, string(bytes))
	}

	if r.Method == "POST" {
		var record Record
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&record); err != nil {
			handleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		record.User_id = auth.Api_user
		ID, err := writeRecord(dbCon, record)
		if err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		record.ID = int(ID)
		respondWithJSON(w, r, http.StatusCreated, record)
	}
}

func handleSingleRecord(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["id"])
	if err != nil {
		handleError(w, r, http.StatusNotFound, err.Error())
		return
	}

	if r.Method == "GET" {
		record, err := getRecord(dbCon, ID)
		if err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		bytes, err := json.MarshalIndent(record, "", "  ")
		if err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		io.WriteString(w, string(bytes))
	}

	if r.Method == "DELETE" {
		if err := deleteRecord(dbCon, ID); err != nil {
			handleError(w, r, http.StatusNotFound, err.Error())
			return
		}
		respondWithJSON(w, r, http.StatusNoContent, ID)
	}

	if r.Method == "PUT" {
		var record Record
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&record); err != nil {
			handleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		record.ID = ID
		if err := updateRecord(dbCon, record); err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, r, http.StatusOK, record)
	}
}
