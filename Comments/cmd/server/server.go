package main

import (
	"Comments/pkg/api"
	"Comments/pkg/db"
	"log"
	"net/http"
	"os"
)

type server struct {
	db  *db.Storage
	api *api.API
}

func main() {
	file, err := os.OpenFile("comments.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Не получилось открыть log файл:", err)
	}
	log.SetOutput(file)
	var srv server
	db1, err := db.New("postgres://postgres:password@192.168.1.191:5432/Comments")
	if err != nil {
		log.Fatal(err)
	}
	srv.db = db1
	srv.api = api.New(srv.db)
	err = http.ListenAndServe(":8082", srv.api.Router())
	log.Println("Запуск сервера на порту :8082")
	if err != nil {
		log.Fatal(err)
	}
}
