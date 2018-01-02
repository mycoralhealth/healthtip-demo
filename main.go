package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func main() {

	dbPath := os.Getenv("DB")
	fmt.Println("opening database: ", dbPath)

	dbCon, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dbCon.Close()

	log.Fatal(run(dbCon))

}
