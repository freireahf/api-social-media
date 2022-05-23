package db

import (
	"api/src/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //Driver
)

//CreateConnection open connection with database mysql
func CreateConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.Connection)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
