package record

import (
	"database/sql"
	"time"
	_ "github.com/lib/pq"
)

type Record struct {
	Id 			int			`json:"id"`
	Username 	string		`json:"username"`
	BestScore 	int			`json:"bestScore"`
	CreatedAt  	time.Time	`json:"createdAt"`
	UpdatedAt  	time.Time	`json:"updatedAt"`
}

type Result struct {
	Data 		interface{}
	StatusCode 	int
}

func InsertRow(db *sql.DB, record Record) Result {
	var result Result
	// Check if the username is unique and if the new Record is valid.
	if checkUsernameExists(db, record.Username) || 
	record.Username == "" || record.BestScore <= 0 {
		result.StatusCode = 400
		return result
	}
	// Insert a row into the "tbl_record" table.
	if _, err := db.Exec(
		`INSERT INTO tbl_record (username, best_score, created_at, updated_at) 
		VALUES ($1, $2, NOW(), NOW())`, record.Username, record.BestScore); err != nil {
		result.StatusCode = 500
		return result
	}
	result.StatusCode = 201
	return result
}

func checkUsernameExists(db *sql.DB, username string) bool {
	row := db.QueryRow(`SELECT count(id) as usernameUsed FROM tbl_record WHERE username=$1`, username)
	var usernameUsed int
	switch err := row.Scan(&usernameUsed); err {
		case sql.ErrNoRows:
			return false
		case nil:
			if usernameUsed > 0 {
				return true
			}
			return false
		default: // TODO
			panic(err)
	}
}

func SelectAllRows(db *sql.DB) Result {
	var result Result
	rows, err := db.Query("SELECT id, username, best_score, created_at, updated_at FROM tbl_record order by best_score desc")
	if err != nil {
		result.StatusCode = 500
		return result
	}
	defer rows.Close()
	records := make([]Record, 0)
	for rows.Next() {
		var record Record
		if err := rows.Scan(&record.Id, &record.Username, &record.BestScore,
							&record.CreatedAt, &record.UpdatedAt); err != nil {
			result.StatusCode = 500
			return result
		}
		records = append(records, record)
	}
	result.StatusCode = 200
	result.Data = records
	return result
}

func SelectOneRow(db *sql.DB, username string) Result {
	var result Result
	row := db.QueryRow(`SELECT id, username, best_score, created_at, updated_at FROM tbl_record WHERE username=$1`, username)
	var record Record
	switch err := row.Scan(&record.Id, &record.Username, &record.BestScore,
		&record.CreatedAt, &record.UpdatedAt); err {
		case sql.ErrNoRows:
			result.StatusCode = 404
		case nil:
			result.StatusCode = 200
			result.Data = record
		default:
			result.StatusCode = 500
	}
	return result
}

func UpdateRow(db *sql.DB, username string, record Record) Result {
	var result Result

	row := db.QueryRow(`SELECT username, best_score FROM tbl_record WHERE username=$1`, username)
	var userRecord Record
	row.Scan(&userRecord.Username, &userRecord.BestScore);
	// Check if the new score is valid and if the username exists
	if record.BestScore <= userRecord.BestScore || userRecord.Username == "" {
		result.StatusCode = 400
		return result
	}
	if _, err := db.Exec(`UPDATE tbl_record SET 
	best_score = $1, updated_at = NOW() WHERE username = $2`, 
	record.BestScore, username); err != nil {
		result.StatusCode = 500
		return result
	}
	result.StatusCode = 200
	return result
}

func DeleteRow(db *sql.DB, username string) Result {
	var result Result
	if checkUsernameExists(db, username) == false {
		result.StatusCode = 400
		return result
	}
	if _, err := db.Exec(`DELETE from tbl_record 
	WHERE username = $1`, username); err != nil {
		result.StatusCode = 500
		return result
	}
	result.StatusCode = 200
	return result
}
