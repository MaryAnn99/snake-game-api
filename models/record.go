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

func InsertRow(db *sql.DB, record Record) error {
	// Insert a record into the "tbl_record" table.
	_, err := db.Exec(
		`INSERT INTO tbl_record (username, best_score, created_at, updated_at) 
		VALUES ($1, $2, NOW(), NOW())`, record.Username, record.BestScore)
	return err
}

func SelectAllRows(db *sql.DB) ([]Record, error) {
	records := make([]Record, 0)
	rows, err := db.Query("SELECT id, username, best_score, created_at, updated_at FROM tbl_record order by best_score desc")
	if err != nil {
		return records, err
	}
	defer rows.Close()
	for rows.Next() {
		var record Record
		if err := rows.Scan(&record.Id, &record.Username, &record.BestScore,
							&record.CreatedAt, &record.UpdatedAt); err != nil {
			return records, err
		}
		records = append(records, record)
	}
	return records, nil
}

func SelectOneRow(db *sql.DB, username string) (Record, error) {
	row := db.QueryRow(`SELECT id, username, best_score, created_at, updated_at FROM tbl_record WHERE username=$1`, username)
	var record Record
	err := row.Scan(&record.Id, &record.Username, &record.BestScore,
		&record.CreatedAt, &record.UpdatedAt)
	return record, err
}

func UpdateRow(db *sql.DB, username string, record Record) error {
	_, err := db.Exec(`UPDATE tbl_record SET 
	best_score = $1, updated_at = NOW() WHERE username = $2`, 
	record.BestScore, username)
	return err
}

func DeleteRow(db *sql.DB, username string) error {
	_, err := db.Exec(`DELETE from tbl_record 
	WHERE username = $1`, username)
	return err
}
