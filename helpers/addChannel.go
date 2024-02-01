package helpers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func AddChannel(writer http.ResponseWriter, request *http.Request) (string, error) {

	channel := request.URL.Query().Get("channel")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return "", err
	}
	defer db.Close()

	randomToken, err := generateRandomToken()
	if err != nil {
		fmt.Println("Error generating random token:", err)
		return "", err
	}

	_, err = db.Exec("INSERT OR REPLACE INTO tokens (channel, token) VALUES (?, ?)", channel, randomToken)
	if err != nil {
		fmt.Println("Error storing token in database:", err)
		return "", err
	}

	return randomToken, nil
}

func generateRandomToken() (string, error) {

	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(randomBytes)

	return token, nil
}
