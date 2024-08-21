package main

import (
	"API_Gateway/pkg/api"
	"log"
	"net/http"
	"os"
)

type server struct {
	api *api.API
}

func main() {
	file, err := os.OpenFile("apigateway.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Не получилось открыть log файл:", err)
	}
	log.SetOutput(file)
	var srv server
	srv.api = api.New()
	http.ListenAndServe(":8080", srv.api.Router())
}
