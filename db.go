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

func hashPassword(password string) string {
	hash := sha256.New()
	saltedPassword := "$%&*)(@#$)(*%@" + password + "%#$(*&#$%(*&@#)%"
	hash.Write([]byte(saltedPassword))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}
