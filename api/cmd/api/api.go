package main

import (
	"fmt"
	"github.com/ivorscott/go-delve-reload/internal/models/postgres"
	"github.com/ivorscott/go-delve-reload/pkg/secrets"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	products *postgres.ProductModel
}

func main() {
	dockerSecrets, _ := secrets.NewDockerSecrets()

	dbname, _ := dockerSecrets.Get("postgres_db")
	dbuser, _ := dockerSecrets.Get("postgres_user")
	dbpass, _ := dockerSecrets.Get("postgres_passwd")
	dbhost, _ := dockerSecrets.Get("postgres_host")

	addr := ":" + os.Getenv("ADDR_PORT")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", dbuser, dbpass, dbhost, dbname)
	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		errorLog.Panicf("[ Database Not Connected ] %s", err)
	}
	defer db.Close()

	err = db.DB().Ping()
	if err != nil {
		errorLog.Panicf("[ Database Not Connected ] %s", err)
	} else {
		infoLog.Printf("[ Database Connected ]")
	}

	dbase := db.DB()
	defer dbase.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		products: &postgres.ProductModel{DB: db},
	}
	server := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("[ Starting server on %s ]", addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}
