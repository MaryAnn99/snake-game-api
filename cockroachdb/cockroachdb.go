package cockroachdb

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	driverName := "postgres"
	dataSourceName  := "postgresql://root@localhost:26257/snake_game_db?sslmode=disable"
	
	// Connect to the "snake_game_db" database.
	db, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		// panic(err)
		db = nil
	}
	if err = db.Ping(); err != nil {
		// panic(err)
		db = nil
	}
}

func GetDatabaseConnection() *sql.DB {
	return db
}