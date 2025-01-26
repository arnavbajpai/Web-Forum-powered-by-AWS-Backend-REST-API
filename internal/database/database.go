package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DBCon *sql.DB
)

func InitializeDB() error {
	var err error
	DBCon, err = sql.Open("mysql", "root:<Password>@tcp(127.0.0.1:3306)/web_forum?parseTime=true")
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	if err = DBCon.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}
