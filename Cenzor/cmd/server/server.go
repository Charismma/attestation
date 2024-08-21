package main

import (
	"Cenzor/pkg/api"
	"log"
	"net/http"
	"os"
)

type server struct {
	api *api.API
}

func main() {
	file, err := os.OpenFile("cenzor.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Не получилось открыть log файл:", err)
	}
	log.SetOutput(file)
	var srv server
	srv.api = api.New()
	err = http.ListenAndServe(":8083", srv.api.Router())
	log.Println("Запуск сервера на порту :8083")
	if err != nil {
		log.Fatal(err)
	}
}
