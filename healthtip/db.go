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

	result, err := dbCon.Exec(`INSERT INTO users (email, first_name, last_name, password, last_tip_epoch) VALUES ($1, $2, $3, $4, 0);`, u.Email, u.FirstName, u.LastName, hashPassword(u.Password))
	if err != nil {
		return 0, fmt.Errorf("couldn't insert %v into users table: %v", u.Email, err)
	}

	Id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("couldn't retrieve Id of insert user %v : %v", u.Email, err)
	}

	return Id, nil
}

func writeAuthToken(dbCon *sql.DB, auth AuthToken) error {
	_, err := dbCon.Exec(`INSERT INTO auth_tokens (api_user, api_key) VALUES ($1, $2);`, auth.ApiUser, auth.ApiKey)
	if err != nil {
		return fmt.Errorf("couldn't insert auth token for user: %v inth auth table", auth.ApiUser)
	}

	return nil
}

func checkAPIAuth(dbCon *sql.DB, auth AuthToken) error {

	var a AuthToken
	if err := dbCon.QueryRow(`SELECT * FROM auth_tokens WHERE api_user = $1 AND api_key = $2;`, auth.ApiUser, auth.ApiKey).Scan(&a); err == sql.ErrNoRows {
		return fmt.Errorf("Incorrect API token for user: %v", auth.ApiUser)
	}
	return nil
}

// returnAuthToken fills out the Auth_user if only the Auth_key is available
func returnAuthUserId(dbCon *sql.DB, auth AuthToken) (int, error) {

	var a AuthToken
	if err := dbCon.QueryRow(`SELECT * FROM auth_tokens WHERE api_key = $1;`, auth.ApiKey).Scan(&a.ApiUser, &a.ApiKey); err == sql.ErrNoRows {
		return 0, fmt.Errorf("Invalid token")
	}
	return a.ApiUser, nil
}

func deleteAuthToken(dbCon *sql.DB, auth AuthToken) error {
	_, err := dbCon.Exec(`DELETE FROM auth_tokens WHERE Api_user = $1;`, auth.ApiUser)

	if err != nil {
		return fmt.Errorf("couldn't delete user %v in auth_tokens table", auth.ApiUser)
	}

	return nil
}

func updateUser(dbCon *sql.DB, u User) error {

	_, err := dbCon.Exec(`UPDATE users SET first_name = $1, last_name = $2, password = $3 WHERE ROWID = $4 ;`, u.FirstName, u.LastName, hashPassword(u.Password), u.Id)
	if err != nil {
		return fmt.Errorf("couldn't update %v in users table: %v", u.Email, err)
	}

	return nil
}

func updateUserTipTime(dbCon *sql.DB, u User) error {

	_, err := dbCon.Exec(`UPDATE users SET last_tip_epoch = $1 WHERE ROWID = $2 ;`, u.LastTip, u.Id)
	if err != nil {
		return fmt.Errorf("couldn't update %v in users table: %v", u.Email, err)
	}

	return nil
}

func updateUserPassword(dbCon *sql.DB, u User) error {

	_, err := dbCon.Exec(`UPDATE users SET password = $1 WHERE ROWID = $2 ;`, hashPassword(u.Password), u.Id)
	if err != nil {
		return fmt.Errorf("couldn't update %v in users table: %v", u.Email, err)
	}

	return nil
}

func checkUserExists(dbCon *sql.DB, u User) (User, error) {
	var user User
	if err := dbCon.QueryRow(`SELECT ROWID, email, first_name, last_name, password FROM users WHERE email = $1;`, strings.ToLower(u.Email)).Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Password); err == sql.ErrNoRows {
		return user, fmt.Errorf("The account doesn't exit: %v", u.Email)
	}

	return user, nil
}

func findUser(dbCon *sql.DB, Id int) (User, error) {
	var user User
	if err := dbCon.QueryRow(`SELECT ROWID, email, first_name, last_name FROM users WHERE ROWID = $1;`, Id).Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName); err == sql.ErrNoRows {
		return user, fmt.Errorf("user not found: %v", Id)
	}

	return user, nil
}

func getUserForId(dbCon *sql.DB, userId int) (User, error) {
	var user User
	if err := dbCon.QueryRow(`SELECT ROWID, * FROM users WHERE ROWID = $1;`, userId).Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.LastTip); err == sql.ErrNoRows {
		return user, fmt.Errorf("user not found: %v", userId)
	}

	return user, nil
}

func checkLoginAuth(dbCon *sql.DB, u User) error {
	var usr User
	if err := dbCon.QueryRow(`SELECT * FROM users WHERE email = $1 AND password = $2;`, strings.ToLower(u.Email), hashPassword(u.Password)).Scan(&usr); err == sql.ErrNoRows {
		return fmt.Errorf("user not found: %v", u.Id)
	}
	return nil
}

func getAllRecords(userId int, dbCon *sql.DB) ([]Record, error) {
	records := make([]Record, 0)
	rows, err := dbCon.Query(`SELECT ROWID, * FROM records WHERE User_id= $1;`, userId)
	if err != nil {
		return records, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.Id, &r.UserId, &r.Age, &r.Height, &r.Weight, &r.Cholesterol, &r.BloodPressure, &r.TipSent, &r.NumberOfCysts, &r.Baldness, &r.BaldnessFromDisease); err != nil {
			return nil, err
		}

		records = append(records, r)
	}

	return records, nil

}

func getRecord(dbCon *sql.DB, Id int) (Record, error) {

	var record Record
	row := dbCon.QueryRow(`SELECT ROWID, * FROM records WHERE ROWID = $1;`, Id)

	if err := row.Scan(&record.Id, &record.UserId, &record.Age, &record.Height, &record.Weight, &record.Cholesterol, &record.BloodPressure, &record.TipSent, &record.NumberOfCysts, &record.Baldness, &record.BaldnessFromDisease); err != nil {
		return record, err
	}

	return record, nil

}

func deleteRecord(dbCon *sql.DB, Id int) error {

	if err := checkRecordExists(dbCon, Id); err != nil {
		return err
	}

	_, err := dbCon.Exec(`DELETE FROM records WHERE ROWID = $1;`, Id)

	if err != nil {
		return fmt.Errorf("couldn't delete %v in records table: %v", Id, err)
	}

	return nil

}

func updateRecord(dbCon *sql.DB, record Record) error {

	if err := checkRecordExists(dbCon, record.UserId); err != nil {
		return err
	}

	_, err := dbCon.Exec(`UPDATE records
		SET age = $1, height = $2, weight = $3, cholesterol = $4, blood_pressure = $5,
		tip_sent = $6, number_of_cysts = $7, baldness = $8, baldness_from_disease = $9 WHERE ROWID = $10;`,
		record.Age, record.Height, record.Weight, record.Cholesterol, record.BloodPressure, record.TipSent,
		record.NumberOfCysts, record.NumberOfCysts, record.Baldness, record.BaldnessFromDisease, record.Id)

	if err != nil {
		return fmt.Errorf("couldn't record %v : %v", record, err)
	}

	return nil

}

func writeRecord(dbCon *sql.DB, record Record) (int64, error) {
	result, err := dbCon.Exec(`INSERT INTO records
		(user_id, age, height, weight, cholesterol, blood_pressure, tip_sent, number_of_cysts, baldness, baldness_from_disease) VALUES ($1, $2, $3, $4, $5, $6, 0, $7, $8, $9);
	`, record.UserId, record.Age, record.Height, record.Weight, record.Cholesterol, record.BloodPressure, record.NumberOfCysts, record.Baldness, record.BaldnessFromDisease)
	if err != nil {
		return 0, fmt.Errorf("couldn't insert %v into records table: %v", record, err)
	}
	Id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("couldn't retrieve ID of inserted record %v: %v", record, err)
	}

	return Id, nil
}

func checkRecordExists(dbCon *sql.DB, Id int) error {
	var record Record

	if err := dbCon.QueryRow(`SELECT id FROM records WHERE ROWID = $1;
	`, Id).Scan(&record); err == sql.ErrNoRows {
		return fmt.Errorf("record ID %d not found", Id)
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
