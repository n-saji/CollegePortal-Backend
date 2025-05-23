package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Postgres_User         string
	Postgres_Db_Name      string
	Postgres_Password     string
	Postgres_Port         int64
	Postgres_Host         string
	Port                  string
	DB_URL                string
	AccountTypeInstructor = "instructor"
	AccountTypeStudent    = "student"
	FRONTEND_URL          string
)

func Init() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file - ", err)
	}

	Postgres_User = os.Getenv("POSTGRES_USER")
	Postgres_Db_Name = os.Getenv("POSTGRES_DB_NAME")
	Postgres_Password = os.Getenv("POSTGRES_PASSWORD")
	Postgres_Port, _ = strconv.ParseInt(os.Getenv("POSTGRES_PORT"), 10, 64)
	Postgres_Host = os.Getenv("POSTGRES_HOST")
	Port = ":" + os.Getenv("PORT")
	SSL_MODE := os.Getenv("SSL_MODE")
	FRONTEND_URL = os.Getenv("FRONTEND_URL")

	os.Setenv("db_url", fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", Postgres_Host, Postgres_User, Postgres_Password, Postgres_Db_Name, Postgres_Port, SSL_MODE))
	DB_URL = os.Getenv("db_url")
}
