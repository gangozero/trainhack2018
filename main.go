package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
)

//go:generate esc -o static.go -pkg main static

func main() {
	log.Println("App started")

	var dbConf dbConfig
	err := envconfig.Process("db", &dbConf)
	if err != nil {
		log.Fatalf("Cannot read DB configuration: %s", err.Error())
	}

	pool, err := newDbConn(dbConf)
	if err != nil {
		log.Fatalf("Cannot start DB connection: %s", err.Error())
	}

	err = dbCheck(pool)
	if err != nil {
		log.Fatalf("Cannot check DB connection: %s", err.Error())
	}

	log.Println("DB connected")
	s := newServer(pool)

	http.HandleFunc("/list", s.handelList())
	http.HandleFunc("/order", s.handelOrder())
	http.HandleFunc("/tasks", s.handelTasks())
	http.Handle("/static/", http.FileServer(FS(false)))
	log.Printf("Starting server on port %s", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
