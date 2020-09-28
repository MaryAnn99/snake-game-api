package cockroachdb

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func init() {
	DB = connectToDB()
}

func connectToDB() *sql.DB {
	driverName := "postgres"
	dataSourceName  := "postgresql://root@localhost:26257/snake_game_db?sslmode=disable"
	// Connect to the "snake_game_db" database.
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	return db
}