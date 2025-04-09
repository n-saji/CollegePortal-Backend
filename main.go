package main

import (
	"CollegeAdministration/config"
	"CollegeAdministration/handlers"
	"CollegeAdministration/jobs"
	"CollegeAdministration/utils"
	"database/sql"
	"embed"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gopkg.in/robfig/cron.v2"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {

	config.Init()
	db := config.DBInit()
	toRunGooseMigration(config.DB_URL)

	handlerConnection := handlers.New(db)

	s := cron.New()
	go s.AddFunc("@every 24h", jobs.RunDailyMigrations)
	// go jobs.AccountDetailsMigration(db) //Fix the migration
	go utils.InitiateWebSockets()
	go s.AddFunc("@every 10s", jobs.SendMessages)
	s.Start()

	r := handlerConnection.GetRouter()
	log.Println(config.Port)
	main_err := r.Run(config.Port)
	if main_err != nil {
		log.Println("MAIN - ERROR ", main_err)
		return
	}
	s.Stop()
	defer config.CloseDB(db)

}

func toRunGooseMigration(url string) {

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Println("db conn failed")
		panic(err)
	}
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Println("Setting Goose Postgres Dialect Failed")
		panic(err)
	}
	if err := goose.Up(db, "migrations", goose.WithAllowMissing()); err != nil {
		log.Println("Goose Up Failed")
		panic(err)
	}
	db.Close()
}
