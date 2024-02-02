package db

import (
	"fmt"
    "log"
	"github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

var CONNECTION *sqlx.DB

func TestConnection() error {
    return CONNECTION.Ping()
}

func ConnectDB(host string, port int, user string, password string, dbname string) *sqlx.DB {

	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf(
            "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", 
            host,
            port, 
            user, 
            password, 
            dbname,
        ),
	)

	if err != nil {
        log.Fatalln("Error connecting to database: ", err)
	}

	return db
}

func InitDB(db *sqlx.DB) error {

	_, err := db.Exec(DB_SCHEMA)
	if err != nil {
		return err
	}

	return nil
}
