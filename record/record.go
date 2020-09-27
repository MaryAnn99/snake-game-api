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

func InsertRow(db *sql.DB, record Record) bool {
	// Check if the username is unique and if the new Record is valid.
	if checkUsernameExists(db, record.Username) || 
	record.Username == "" || record.BestScore <= 0 {
		return false
	}
	// Insert a row into the "tbl_record" table.
	if _, err := db.Exec(
		`INSERT INTO tbl_record (username, best_score, created_at, updated_at) 
		VALUES ($1, $2, NOW(), NOW())`, record.Username, record.BestScore); err != nil {
		panic(err)
	}
	return true
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
		default:
			panic(err)
	}
}

func SelectAllRows(db *sql.DB) []Record {
	rows, err := db.Query("SELECT id, username, best_score, created_at, updated_at FROM tbl_record order by best_score desc")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	results := make([]Record, 0)
	for rows.Next() {
		var result Record
		if err := rows.Scan(&result.Id, &result.Username, &result.BestScore,
							&result.CreatedAt, &result.UpdatedAt); err != nil {
			panic(err)
		}
		results = append(results, result)
	}
	return results
}

func SelectOneRow(db *sql.DB, username string) Record {
	row := db.QueryRow(`SELECT id, username, best_score, created_at, updated_at FROM tbl_record WHERE username=$1`, username)
	var record Record
	switch err := row.Scan(&record.Id, &record.Username, &record.BestScore,
		&record.CreatedAt, &record.UpdatedAt); err {
		case sql.ErrNoRows:
			return record
		case nil:
			return record
		default:
			panic(err)
	}
}

func UpdateRow(db *sql.DB, username string, record Record) bool {
	userRecord := SelectOneRow(db, username)
	// Check if the new score is valid and if the username exists
	if record.BestScore <= userRecord.BestScore || userRecord.Username == "" {
		return false
	}
	if _, err := db.Exec(`UPDATE tbl_record SET 
	best_score = $1, updated_at = NOW() WHERE username = $2`, 
	record.BestScore, username); err != nil {
		panic(err)
	}
	return true
}

func DeleteRow(db *sql.DB, username string) bool {
	if checkUsernameExists(db, username) == false {
		return false
	}
	if _, err := db.Exec(`DELETE from tbl_record 
	WHERE username = $1`, username); err != nil {
		panic(err)
	}
	return true
}
