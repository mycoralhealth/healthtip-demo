package healthtip

import (
	"database/sql"
	"fmt"
)

func getUserTip(dbCon *sql.DB, userId string) (int64, error) {
	var timestamp int64
	if err := dbCon.QueryRow(`SELECT timestamp FROM tips WHERE user_id = $1;`, userId).Scan(&timestamp); err == sql.ErrNoRows {
		return 0, fmt.Errorf("User id %v has no tips.", userId)
	}
	return timestamp, nil
}

func updateUserTipTime(dbCon *sql.DB, userId string, lastTip, now int64) error {
	if lastTip != 0 {
		if _, err := dbCon.Exec(`INSERT INTO tips (user_id, timestamp) VALUES ($1, $2);`, userId, now); err != nil {
			return fmt.Errorf("Could not insert tip record for user %v: %v", userId, err)
		}
	} else {
		if _, err := dbCon.Exec(`UPDATE tips SET timestamp = $1 WHERE user_id = $2;`, now, userId); err != nil {
			return fmt.Errorf("Could not update tip record for user %v: %v", userId, err)
		}
	}
	return nil
}

func getAllRecords(userId string, dbCon *sql.DB) ([]Record, error) {
	records := make([]Record, 0)
	rows, err := dbCon.Query(`SELECT ROWID, * FROM records WHERE user_id= $1;`, userId)
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

	if err := checkRecordExists(dbCon, record.Id); err != nil {
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
