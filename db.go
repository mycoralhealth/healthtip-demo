package main

import (
	"database/sql"
	"fmt"
)

func writeUser(dbCon *sql.DB, u User) (int64, error) {

	result, err := dbCon.Exec(`INSERT INTO users (email, first_name, last_name, password) VALUES ($1, $2, $3, $4);`, u.Email, u.First_name, u.Last_name, u.Password)
	if err != nil {
		return 0, fmt.Errorf("couldn't insert %v into users table: %v", u.Email, err)
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("couldn't retrieve ID of insert user %v : %v", u.Email, err)
	}

	return ID, nil
}
