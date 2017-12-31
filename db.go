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

func updateUser(dbCon *sql.DB, u User) error {

	_, err := dbCon.Exec(`UPDATE users SET first_name = $1, last_name = $2, password = $3;`, u.First_name, u.Last_name, hashPassword(u.Password))
	if err != nil {
		return fmt.Errorf("couldn't update %v in users table: %v", u.Email, err)
	}

	return nil
}

func checkUserExists(dbCon *sql.DB, u User) error {
	var user User
	if err := dbCon.QueryRow(`SELECT ROWID FROM users WHERE email = $1;`, u.Email).Scan(&user); err == sql.ErrNoRows {
		return fmt.Errorf("user not found: %v", u.Email)
	}
	return nil

}

func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}
