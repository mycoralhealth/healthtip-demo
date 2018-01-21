package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

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

	if len(user.Password) < 6 {
		handleError(w, r, http.StatusUnprocessableEntity, "Minimum password length is 6 characters")
		return
	}

	lastID, err := writeUser(dbCon, user)
	if err != nil {
		handleError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user.ID = int(lastID)
	auth, err := createAuthToken(user.ID, dbCon)
	if err != nil {
		handleError(w, r, http.StatusUnprocessableEntity, err.Error())
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

	auth, _ := getBasicAPIAuth(r)
	user.ID = auth.Api_user

	if err := updateUser(dbCon, user); err != nil {
		handleError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	respondWithJSON(w, r, http.StatusCreated, user.Email)
}

func handleLogin(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {

	user, err := getBasicLoginAuth(r)

	if err != nil {
		handleError(w, r, http.StatusUnauthorized, err.Error())
		return
	}

	dbUser, err := checkUserExists(dbCon, user)

	if err != nil {
		handleError(w, r, 404, err.Error())
		return
	}

	if err := checkLoginAuth(dbCon, user); err != nil {
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

func handleChangePassword(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	auth, _ := getBasicAPIAuth(r)

	var user User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		handleError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	user.ID = auth.Api_user
	err := updateUserPassword(dbCon, user)
	if err != nil {
		handleError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, r, http.StatusNoContent, nil)
}

func handleResetPassword(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	// request should send user object with just email
	var user User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		handleError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	dbUser, err := checkUserExists(dbCon, user)

	if err != nil {
		handleError(w, r, 404, err.Error())
		return
	}

	auth, err := createAuthToken(dbUser.ID, dbCon)
	if err != nil {
		handleError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	url := os.Getenv("CLIENT_URL") + "changePassword?token=" + url.QueryEscape(auth.Api_key)
	emailPasswordReset(dbUser, url)

	respondWithJSON(w, r, http.StatusNoContent, nil)
}

func handleClaimToken(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {

	var auth AuthToken
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&auth); err != nil {
		handleError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	userID, err := returnAuthUserID(dbCon, auth)
	if err != nil {
		handleError(w, r, 404, err.Error())
		return
	}

	auth.Api_user = userID

	if err := deleteAuthToken(dbCon, auth); err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := findUser(dbCon, userID)
	if err != nil {
		handleError(w, r, 404, err.Error())
		return
	}

	newAuth, err := createAuthToken(user.ID, dbCon)
	if err != nil {
		handleError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	respondWithJSON(w, r, http.StatusOK, makeLoginResult(user, newAuth))
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
		return token, errors.New("invalid API Authorization token. Required Basic Api_user:Api_key")
	}

	apiUser, err := strconv.Atoi(p)
	if err != nil {
		return token, errors.New("invalid Api_user")
	}

	token.Api_user = apiUser
	token.Api_key = s

	return token, nil
}

func handleLogout(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	auth, _ := getBasicAPIAuth(r)

	if err := deleteAuthToken(dbCon, auth); err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, r, http.StatusNoContent, nil)
}

func handleRecords(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	auth, _ := getBasicAPIAuth(r)

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

func handleRecordTip(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["id"])

	if err != nil {
		handleError(w, r, http.StatusNotFound, err.Error())
		return
	}

	record, err := getRecord(dbCon, ID)
	if err != nil {
		handleError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if record.Tip_sent != 0 {
		handleError(w, r, http.StatusUnprocessableEntity, "You can only request a tip on the same record once.")
		return
	}

	auth, _ := getBasicAPIAuth(r)

	dbUser, err := getUserForId(dbCon, auth.Api_user)

	if err != nil {
		handleError(w, r, http.StatusNotFound, err.Error())
		return
	}

	now := time.Now()

	if dbUser.Last_tip != 0 {
		secs := now.Unix()
		interval, err := strconv.ParseInt(os.Getenv("TIP_INTERVAL"), 10, 64)

		if err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		if secs < (dbUser.Last_tip + interval) {
			handleError(w, r, http.StatusUnprocessableEntity, "You can only request one Health Tip every 24 hours.")
			return
		}
	}

	mailErr := emailHealthTipRequest(dbUser, record)

	if mailErr != nil {
		handleError(w, r, http.StatusUnprocessableEntity, "Unable to send email. Please contact Coral Health.")
		return
	}

	dbUser.Last_tip = now.Unix()
	updateUserTipTime(dbCon, dbUser)

	record.Tip_sent = 1
	updateRecord(dbCon, record)

	respondWithJSON(w, r, http.StatusOK, record)
}
