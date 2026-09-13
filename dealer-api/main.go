package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"dealer-api/handlers"
	_ "github.com/lib/pq"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	dbHost, dbPort, dbUser, dbPass, dbName)

	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil && db.Ping() == nil {
			log.Println("Successful connect to Postgresql!")
			break
		}
		log.Println("Waiting connecting...")
		time.Sleep(2 * time.Second)
	}

	if err != nil || db.Ping() != nil {
		log.Fatal("Error connect to database ", err)
	}

	authHandler := &handlers.AuthHandler{DB: db}

	http.HandleFunc("/auth/login", authHandler.LoginHandler)
	http.HandleFunc("/auth/register", authHandler.RegisterHandler)

	log.Println("dealer-api запущен на порту 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}