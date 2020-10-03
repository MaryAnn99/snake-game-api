package cockroachdb

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectToDB() {
	driverName := "postgres"
	dataSourceName  := "postgresql://root@localhost:26257/snake_game_db?sslmode=disable"
	// Connect to the "snake_game_db" database.
	var err error
	DB, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	if err = DB.Ping(); err != nil {
		panic(err)
	}
}