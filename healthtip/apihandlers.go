package healthtip

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func makeLoginResult(user User, auth AuthToken) LoginResult {
	var result LoginResult

	result.Token = auth
	result.Email = user.Email
	result.FirstName = user.FirstName
	result.LastName = user.LastName

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

	lastId, err := writeUser(dbCon, user)
	if err != nil {
		handleError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user.Id = int(lastId)
	auth, err := createAuthToken(user.Id, dbCon)
	if err != nil {
		handleError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	respondWithJSON(w, r, http.StatusCreated, makeLoginResult(user, auth))
}

func createAuthToken(Id int, dbCon *sql.DB) (AuthToken, error) {
	var auth AuthToken
	// create auth token and put in DB
	auth.ApiUser = Id
	// generated salted hash from current time and random number is the auth token
	token, err := getRandomString(10)
	if err != nil {
		return auth, err
	}

	auth.ApiKey = hashPassword(token)
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
	user.Id = auth.ApiUser

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

	auth, err := createAuthToken(dbUser.Id, dbCon)
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

	user.Id = auth.ApiUser
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

	auth, err := createAuthToken(dbUser.Id, dbCon)
	if err != nil {
		handleError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	url := os.Getenv("CLIENT_URL") + "changePass?token=" + url.QueryEscape(auth.ApiKey)
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

	userId, err := returnAuthUserId(dbCon, auth)
	if err != nil {
		handleError(w, r, 404, err.Error())
		return
	}

	auth.ApiUser = userId

	if err := deleteAuthToken(dbCon, auth); err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := findUser(dbCon, userId)
	if err != nil {
		handleError(w, r, 404, err.Error())
		return
	}

	newAuth, err := createAuthToken(user.Id, dbCon)
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
		return token, errors.New("invalid API Authorization token. Required Basic ApiUser:ApiKey")
	}

	apiUser, err := strconv.Atoi(p)
	if err != nil {
		return token, errors.New("invalid ApiUser")
	}

	token.ApiUser = apiUser
	token.ApiKey = s

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
		records, err := getAllRecords(auth.ApiUser, dbCon)
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

		record.UserId = auth.ApiUser
		Id, err := writeRecord(dbCon, record)
		if err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		record.Id = int(Id)
		respondWithJSON(w, r, http.StatusCreated, record)
	}
}

func handleSingleRecord(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	vars := mux.Vars(r)
	Id, err := strconv.Atoi(vars["id"])
	if err != nil {
		handleError(w, r, http.StatusNotFound, err.Error())
		return
	}

	if r.Method == "GET" {
		record, err := getRecord(dbCon, Id)
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
		if err := deleteRecord(dbCon, Id); err != nil {
			handleError(w, r, http.StatusNotFound, err.Error())
			return
		}
		respondWithJSON(w, r, http.StatusNoContent, Id)
	}

	if r.Method == "PUT" {
		var record Record
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&record); err != nil {
			handleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		record.Id = Id
		if err := updateRecord(dbCon, record); err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, r, http.StatusOK, record)
	}
}

func handleRecordTip(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	vars := mux.Vars(r)
	Id, err := strconv.Atoi(vars["id"])

	if err != nil {
		handleError(w, r, http.StatusNotFound, err.Error())
		return
	}

	record, err := getRecord(dbCon, Id)
	if err != nil {
		handleError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if record.TipSent != 0 {
		handleError(w, r, http.StatusUnprocessableEntity, "You can only request a tip on the same record once.")
		return
	}

	auth, _ := getBasicAPIAuth(r)

	dbUser, err := getUserForId(dbCon, auth.ApiUser)

	if err != nil {
		handleError(w, r, http.StatusNotFound, err.Error())
		return
	}

	now := time.Now()

	if dbUser.LastTip != 0 {
		secs := now.Unix()
		interval, err := strconv.ParseInt(os.Getenv("TIP_INTERVAL"), 10, 64)

		if err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		if secs < (dbUser.LastTip + interval) {
			handleError(w, r, http.StatusUnprocessableEntity, "You can only request one Health Tip every 24 hours.")
			return
		}
	}

	mailErr := emailHealthTipRequest(dbUser, record)

	if mailErr != nil {
		handleError(w, r, http.StatusUnprocessableEntity, "Unable to send email. Please contact Coral Health.")
		return
	}

	dbUser.LastTip = now.Unix()
	updateUserTipTime(dbCon, dbUser)

	record.TipSent = 1
	updateRecord(dbCon, record)

	respondWithJSON(w, r, http.StatusOK, record)
}

const (
	// Supported treatments
	HairRemoval    int = 1
	HairTransplant int = 2
)

func handleCompanies(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	if r.Method == "GET" {
		companies, err := getAllCompanies(dbCon)
		if err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, r, http.StatusOK, companies)
	}
}

func handleProcedures(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	if r.Method == "GET" {
		procedures, err := getAllProcedures(dbCon)
		if err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, r, http.StatusOK, procedures)
	}
}

func handleInsuranceApproval(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	vars := mux.Vars(r)
	Id, err := strconv.Atoi(vars["id"])

	if err != nil {
		handleError(w, r, http.StatusNotFound, err.Error())
		return
	}

	record, err := getRecord(dbCon, Id)
	if err != nil {
		handleError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if r.Method == "POST" {
		var request InsuranceApprovalRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&request); err != nil {
			handleError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		response := InsuranceApprovalResponse{InsuranceApprovalRequest: request}
		var approved bool
		var err error
		if approved, err = getConditionalApproval(dbCon, request.Company.Id, request.Procedure.Id); err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		// Approval logic
		switch request.Procedure.Id {
		case HairRemoval:
			response.Approved = (record.NumberOfCysts > 1 && approved)
		case HairTransplant:
			response.Approved = (record.BaldnessFromDisease && approved)
		default:
			handleError(w, r, http.StatusBadRequest, fmt.Sprintf("Unsupported Procedure: value=%+v", request))
			return
		}
		respondWithJSON(w, r, http.StatusOK, response)
	}

}

func handleCompanyPolicy(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	companyId, err := strconv.Atoi(vars["companyId"])
	if err != nil {
		handleError(w, r, http.StatusNotFound, err.Error())
		return
	}
	procedureId, err := strconv.Atoi(vars["procedureId"])
	if err != nil {
		handleError(w, r, http.StatusNotFound, err.Error())
		return
	}

	// Get the file bytes
	name, fileBytes, err := getPolicyFile(db, companyId, procedureId)
	if err != nil {
		handleError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// write file to download stream
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	w.Header().Set("Content-Disposition", "attachment; filename="+name)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(w, bytes.NewReader(fileBytes))
	if err != nil {
		handleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

}