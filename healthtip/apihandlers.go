package healthtip

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func getUserId(r *http.Request) string {
	return r.Context().Value("userId").(string)
}

func handleRecords(w http.ResponseWriter, r *http.Request, dbCon *sql.DB) {
	userId := getUserId(r)

	if r.Method == "GET" {
		records, err := getAllRecords(userId, dbCon)
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

		record.UserId = userId
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

	userId := getUserId(r)
	lastTip, _ := getUserTip(dbCon, userId)
	now := time.Now().Unix()

	if lastTip != 0 {
		interval, err := strconv.ParseInt(os.Getenv("TIP_INTERVAL"), 10, 64)

		if err != nil {
			handleError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		if now < (lastTip + interval) {
			handleError(w, r, http.StatusUnprocessableEntity, "You can only request one Health Tip every 24 hours.")
			return
		}
	}

	mailErr := emailHealthTipRequest(userId, record)

	if mailErr != nil {
		handleError(w, r, http.StatusUnprocessableEntity, "Unable to send email. Please contact Coral Health.")
		return
	}

	updateUserTipTime(dbCon, userId, lastTip, now)

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
