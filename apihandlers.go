package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

	respondWithJSON(w, r, http.StatusCreated, user.Email)

}

func handleLogin(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	err := r.ParseForm()
	if err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return

	}

	//	fmt.Println(ioutil.ReadAll(r.Body))
	var user User
	user.Email = r.FormValue("email")
	user.Password = hashPassword(r.FormValue("password"))

	if err := checkUserExists(dbCon, user); err != nil {
		handleError(w, r, 404, err.Error())
		return
	}

	// make sure email and password match
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
	if err := checkAuthToken(dbCon, auth); err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if r.Method == "GET" {
		records, err := getAllRecords(dbCon)
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

		ID, err := writeRecord(dbCon, record)
		if err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		record.User_id = int(ID)
		respondWithJSON(w, r, http.StatusCreated, record)

	}

}

func handleSingleRecord(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	if err := checkAuthToken(dbCon, auth); err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

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
			handleError(w, r, 404, err.Error())
			return
		}
		respondWithJSON(w, r, 204, ID)
	}

	if r.Method == "PUT" {
		var record Record
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&record); err != nil {
			handleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		record.User_id = ID
		if err := updateRecord(dbCon, record); err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, r, 200, record)

	}
}
