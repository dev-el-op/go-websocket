package helpers

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "./tokens.db"

func IsValidToken(channel, token string) bool {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Could not open database connection:", err)
		return false
	}
	defer db.Close()

	row := db.QueryRow("SELECT token FROM tokens WHERE channel = ?", channel)

	var storedToken string
	if err := row.Scan(&storedToken); err != nil {
		fmt.Println("Error durgin database query:", err)
		return false
	}

	return token == storedToken
}
