package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
)

func writeUser(dbCon *sql.DB, u User) (int64, error) {

	// check if user exists and return error if it does
	// if checkUserExists function does NOT throw an error it means the user already exists
	if err := checkUserExists(dbCon, u); err == nil {
		return 0, fmt.Errorf("user %v already exists", u.Email)
	}

	result, err := dbCon.Exec(`INSERT INTO users (email, first_name, last_name, password) VALUES ($1, $2, $3, $4);`, u.Email, u.First_name, u.Last_name, hashPassword(u.Password))
	if err != nil {
		return 0, fmt.Errorf("couldn't insert %v into users table: %v", u.Email, err)
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("couldn't retrieve ID of insert user %v : %v", u.Email, err)
	}

	return ID, nil
}

func writeAuthToken(dbCon *sql.DB, auth AuthToken) error {
	_, err := dbCon.Exec(`INSERT INTO auth_tokens (api_user, api_key) VALUES ($1, $2);`, auth.Api_user, auth.Api_key)
	if err != nil {
		return fmt.Errorf("couldn't insert auth token for user: %v inth auth table", auth.Api_user)
	}

	return nil
}

func checkAuthToken(dbCon *sql.DB, auth AuthToken) error {

	var a AuthToken
	if err := dbCon.QueryRow(`SELECT * FROM auth_tokens WHERE api_user = $1 AND api_key = $2;`, auth.Api_user, auth.Api_key).Scan(&a); err == sql.ErrNoRows {
		return fmt.Errorf("user not found: %v", auth.Api_user)
	}
	return nil
}

func deleteAuthToken(dbCon *sql.DB, auth AuthToken) error {
	_, err := dbCon.Exec(`DELETE FROM auth_tokens WHERE Api_user = $1;`, auth.Api_user)

	if err != nil {
		return fmt.Errorf("couldn't delete user %v in auth_tokens table", auth.Api_user)
	}

	return nil
}

func updateUser(dbCon *sql.DB, u User) error {

	_, err := dbCon.Exec(`UPDATE users SET first_name = $1, last_name = $2, password = $3;`, u.First_name, u.Last_name, hashPassword(u.Password))
	if err != nil {
		return fmt.Errorf("couldn't update %v in users table: %v", u.Email, err)
	}

	return nil
}

func checkUserExists(dbCon *sql.DB, u User) error {
	var user User
	if err := dbCon.QueryRow(`SELECT ROWID FROM users WHERE email = $1;`, u.Email).Scan(&user.ID); err == sql.ErrNoRows {
		return fmt.Errorf("user not found: %v", u.Email)
	}
	return nil

}

// checkLogin checks email against password in users table to make sure they match
func checkLogin(dbCon *sql.DB, u User) (int, error) {
	rows, err := dbCon.Query(`SELECT * FROM users WHERE email = $1 AND password = $2;`, u.Email, u.Password)
	if err != nil {
		return 0, fmt.Errorf("user %v password doesn't match", u.Email)
	}
	defer rows.Close()

	var count int
	for rows.Next() {

		// make sure there's only one entry for user
		if count > 1 {
			return 0, fmt.Errorf("more than one entry for user %v", u.Email)
		}

		count++

		var user User
		if err := rows.Scan(&user); err != nil {
			return 0, fmt.Errorf("error scanning rows for user %v", u.Email)
		}

	}

	return u.ID, nil

}

func getAllRecords(dbCon *sql.DB) ([]Record, error) {
	records := make([]Record, 0)
	rows, err := dbCon.Query(`SELECT * FROM records;`)
	if err != nil {
		return records, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Record
		if err := rows.Scan(&r); err != nil {
			return nil, err
		}

		records = append(records, r)
	}

	return records, nil

}

func getRecord(dbCon *sql.DB, ID int) (Record, error) {

	var record Record
	row := dbCon.QueryRow(`SELECT * FROM records WHERE id = $1;`, ID)

	if err := row.Scan(&record); err != nil {
		return record, err
	}

	return record, nil

}

func deleteRecord(dbCon *sql.DB, ID int) error {

	if err := checkRecordExists(dbCon, ID); err != nil {
		return err
	}

	_, err := dbCon.Exec(`DELETE FROM records WHERE id = $1;`, ID)

	if err != nil {
		return fmt.Errorf("couldn't delete %v in records table: %v", ID, err)
	}

	return nil

}

func updateRecord(dbCon *sql.DB, record Record) error {

	if err := checkRecordExists(dbCon, record.User_id); err != nil {
		return err
	}

	_, err := dbCon.Exec(`UPDATE records
		SET age = $1, height = $2, weight = $3, cholesterol = $4, blood_pressure = $5
		WHERE id = $6;
	`, record.Age, record.Height, record.Weight, record.Cholesterol, record.Blood_pressure, record.User_id)

	if err != nil {
		return fmt.Errorf("couldn't record %v : %v", record, err)
	}

	return nil

}

func writeRecord(dbCon *sql.DB, record Record) (int64, error) {
	result, err := dbCon.Exec(`INSERT INTO records
		(age, height, weight, cholesterol, blood_pressure) VALUES ($1, $2, $3, $4, $5);
	`, record.Age, record.Height, record.Weight, record.Cholesterol, record.Blood_pressure)
	if err != nil {
		return 0, fmt.Errorf("couldn't insert %v into records table: %v", record, err)
	}
	ID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("couldn't retrieve ID of inserted record %v: %v", record, err)
	}

	return ID, nil

}

func checkRecordExists(dbCon *sql.DB, ID int) error {
	var record Record

	if err := dbCon.QueryRow(`SELECT id FROM records WHERE id = $1;
	`, ID).Scan(&record); err == sql.ErrNoRows {
		return fmt.Errorf("record ID %d not found", ID)
	}

	return nil

}

func hashPassword(password string) string {
	hash := sha256.New()
	saltedPassword := "$%&*)(@#$)(*%@" + password + "%#$(*&#$%(*&@#)%"
	hash.Write([]byte(saltedPassword))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}
