package healthtip

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
)

func writeUser(dbCon *sql.DB, u User) (int64, error) {

	// check if user exists and return error if it does
	// if checkUserExists function does NOT throw an error it means the user already exists
	if _, err := checkUserExists(dbCon, u); err == nil {
		return 0, fmt.Errorf("user %v already exists", u.Email)
	}

	result, err := dbCon.Exec(`INSERT INTO users (email, first_name, last_name, password, last_tip_epoch) VALUES ($1, $2, $3, $4, 0);`, u.Email, u.First_name, u.Last_name, hashPassword(u.Password))
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

func checkAPIAuth(dbCon *sql.DB, auth AuthToken) error {

	var a AuthToken
	if err := dbCon.QueryRow(`SELECT * FROM auth_tokens WHERE api_user = $1 AND api_key = $2;`, auth.Api_user, auth.Api_key).Scan(&a); err == sql.ErrNoRows {
		return fmt.Errorf("Incorrect API token for user: %v", auth.Api_user)
	}
	return nil
}

// returnAuthToken fills out the Auth_user if only the Auth_key is available
func returnAuthUserID(dbCon *sql.DB, auth AuthToken) (int, error) {

	var a AuthToken
	if err := dbCon.QueryRow(`SELECT * FROM auth_tokens WHERE api_key = $1;`, auth.Api_key).Scan(&a.Api_user, &a.Api_key); err == sql.ErrNoRows {
		return 0, fmt.Errorf("Invalid token")
	}
	return a.Api_user, nil
}

func deleteAuthToken(dbCon *sql.DB, auth AuthToken) error {
	_, err := dbCon.Exec(`DELETE FROM auth_tokens WHERE Api_user = $1;`, auth.Api_user)

	if err != nil {
		return fmt.Errorf("couldn't delete user %v in auth_tokens table", auth.Api_user)
	}

	return nil
}

func updateUser(dbCon *sql.DB, u User) error {

	_, err := dbCon.Exec(`UPDATE users SET first_name = $1, last_name = $2, password = $3 WHERE ROWID = $4 ;`, u.First_name, u.Last_name, hashPassword(u.Password), u.ID)
	if err != nil {
		return fmt.Errorf("couldn't update %v in users table: %v", u.Email, err)
	}

	return nil
}

func updateUserTipTime(dbCon *sql.DB, u User) error {

	_, err := dbCon.Exec(`UPDATE users SET last_tip_epoch = $1 WHERE ROWID = $2 ;`, u.Last_tip, u.ID)
	if err != nil {
		return fmt.Errorf("couldn't update %v in users table: %v", u.Email, err)
	}

	return nil
}

func updateUserPassword(dbCon *sql.DB, u User) error {

	_, err := dbCon.Exec(`UPDATE users SET password = $1 WHERE ROWID = $2 ;`, hashPassword(u.Password), u.ID)
	if err != nil {
		return fmt.Errorf("couldn't update %v in users table: %v", u.Email, err)
	}

	return nil
}

func checkUserExists(dbCon *sql.DB, u User) (User, error) {
	var user User
	if err := dbCon.QueryRow(`SELECT ROWID, email, first_name, last_name, password FROM users WHERE email = $1;`, strings.ToLower(u.Email)).Scan(&user.ID, &user.Email, &user.First_name, &user.Last_name, &user.Password); err == sql.ErrNoRows {
		return user, fmt.Errorf("The account doesn't exit: %v", u.Email)
	}

	return user, nil
}

func findUser(dbCon *sql.DB, ID int) (User, error) {
	var user User
	if err := dbCon.QueryRow(`SELECT ROWID, email, first_name, last_name FROM users WHERE ROWID = $1;`, ID).Scan(&user.ID, &user.Email, &user.First_name, &user.Last_name); err == sql.ErrNoRows {
		return user, fmt.Errorf("user not found: %v", ID)
	}

	return user, nil
}

func getUserForId(dbCon *sql.DB, user_id int) (User, error) {
	var user User
	if err := dbCon.QueryRow(`SELECT ROWID, * FROM users WHERE ROWID = $1;`, user_id).Scan(&user.ID, &user.Email, &user.First_name, &user.Last_name, &user.Password, &user.Last_tip); err == sql.ErrNoRows {
		return user, fmt.Errorf("user not found: %v", user_id)
	}

	return user, nil
}

func checkLoginAuth(dbCon *sql.DB, u User) error {
	var usr User
	if err := dbCon.QueryRow(`SELECT * FROM users WHERE email = $1 AND password = $2;`, strings.ToLower(u.Email), hashPassword(u.Password)).Scan(&usr); err == sql.ErrNoRows {
		return fmt.Errorf("user not found: %v", u.ID)
	}
	return nil
}

func getAllRecords(user_id int, dbCon *sql.DB) ([]Record, error) {
	records := make([]Record, 0)
	rows, err := dbCon.Query(`SELECT ROWID, * FROM records WHERE User_id= $1;`, user_id)
	if err != nil {
		return records, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.ID, &r.User_id, &r.Age, &r.Height, &r.Weight, &r.Cholesterol, &r.Blood_pressure, &r.Tip_sent, &r.Number_of_cysts, &r.Baldness, &r.Baldness_from_disease); err != nil {
			return nil, err
		}

		records = append(records, r)
	}

	return records, nil

}

func getRecord(dbCon *sql.DB, ID int) (Record, error) {

	var record Record
	row := dbCon.QueryRow(`SELECT ROWID, * FROM records WHERE ROWID = $1;`, ID)

	if err := row.Scan(&record.ID, &record.User_id, &record.Age, &record.Height, &record.Weight, &record.Cholesterol, &record.Blood_pressure, &record.Tip_sent, &record.Number_of_cysts, &record.Baldness, &record.Baldness_from_disease); err != nil {
		return record, err
	}

	return record, nil

}

func deleteRecord(dbCon *sql.DB, ID int) error {

	if err := checkRecordExists(dbCon, ID); err != nil {
		return err
	}

	_, err := dbCon.Exec(`DELETE FROM records WHERE ROWID = $1;`, ID)

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
		SET age = $1, height = $2, weight = $3, cholesterol = $4, blood_pressure = $5, tip_sent = $6
		WHERE ROWID = $7;
	`, record.Age, record.Height, record.Weight, record.Cholesterol, record.Blood_pressure, record.Tip_sent, record.ID)

	if err != nil {
		return fmt.Errorf("couldn't record %v : %v", record, err)
	}

	return nil

}

func writeRecord(dbCon *sql.DB, record Record) (int64, error) {
	result, err := dbCon.Exec(`INSERT INTO records
		(user_id, age, height, weight, cholesterol, blood_pressure, tip_sent, number_of_cysts, baldness, baldness_from_disease) VALUES ($1, $2, $3, $4, $5, $6, 0, $7, $8, $9);
	`, record.User_id, record.Age, record.Height, record.Weight, record.Cholesterol, record.Blood_pressure, record.Number_of_cysts, record.Baldness, record.Baldness_from_disease)
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

	if err := dbCon.QueryRow(`SELECT id FROM records WHERE ROWID = $1;
	`, ID).Scan(&record); err == sql.ErrNoRows {
		return fmt.Errorf("record ID %d not found", ID)
	}

	return nil

}

func getAllCompanies(db *sql.DB) ([]InsuranceCompany, error) {
	var companies []InsuranceCompany
	rows, err := db.Query(`SELECT * FROM insurance_companies;`)
	if err != nil {
		return companies, err
	}
	defer rows.Close()

	for rows.Next() {
		var c InsuranceCompany
		if err := rows.Scan(&c.Id, &c.Name); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, nil
}

func getAllProcedures(db *sql.DB) ([]Procedure, error) {
	var procedures []Procedure
	rows, err := db.Query(`SELECT * FROM procedures;`)
	if err != nil {
		return procedures, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Procedure
		if err := rows.Scan(&p.Id, &p.Name); err != nil {
			return nil, err
		}
		procedures = append(procedures, p)
	}
	return procedures, nil
}

func getConditionalApproval(db *sql.DB, companyId, procedureId int) (bool, error) {
	var approved bool
	row := db.QueryRow(`SELECT conditional_approval FROM medical_policies
		WHERE company_id = $1 and procedure_id = $2;`, companyId, procedureId)

	if err := row.Scan(&approved); err != nil {
		return false, err
	}
	return approved, nil
}

func getPolicyFile(db *sql.DB, companyId, procedureId int) (string, []byte, error) {
	var name string
	var bytes []byte
	row := db.QueryRow(`SELECT f.name, f.content FROM medical_policies AS m INNER JOIN files AS f
		ON m.file_id = f.id WHERE m.company_id = $1 AND m.procedure_id = $2;`, companyId, procedureId)

	if err := row.Scan(&name, &bytes); err != nil {
		return "", nil, err
	}
	return name, bytes, nil
}

func hashPassword(password string) string {
	hash := sha256.New()
	saltedPassword := "$%&*)(@#$)(*%@" + password + "%#$(*&#$%(*&@#)%"
	hash.Write([]byte(saltedPassword))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}
