package main

import (
	"GoNews_project/pkg/api"
	"GoNews_project/pkg/db"
	"GoNews_project/pkg/db/memdb"
	"GoNews_project/pkg/db/postgres"
	"GoNews_project/pkg/rss"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// json-конфиг для ссылок с новостями
type jsonConfig struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

// структура нашего сервера
type server struct {
	db  db.Interface
	api *api.API
}

func main() {
	file, err := os.OpenFile("gonews.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Не получилось открыть log файл:", err)
	}
	log.SetOutput(file)
	var srv server
	initStringDb := "postgres://postgres:password@192.168.1.191:5432/GoNews"
	db1, err := postgres.New(initStringDb) //подключаемся к БД Postgres
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("Подключение к БД")
	db2, err := memdb.New() //подключаемся к БД в памяти
	if err != nil {
		log.Fatal(err)
	}
	_, _ = db1, db2
	srv.db = db1
	srv.api = api.New(srv.db)
	var conf jsonConfig
	//log.Println("Чтение файла с сайтами")
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, &conf)
	if err != nil {
		log.Fatal(err)
	}
	chPosts := make(chan []db.Post)
	chError := make(chan error)
	//запуск горутин для парсинга ссылок
	for _, url := range conf.URLS {
		go ParseUrls(url, &srv.db, chPosts, chError, conf.Period)
		//	log.Println("Запуск одной из горутин для парсинга")
	}
	//горутина обработки сообщений из канала с постами
	go func() {
		for posts := range chPosts {
			//log.Println("Добавление постов в базу")
			srv.db.AddPosts(posts)
		}
	}()
	//горутина обработки сообщений из канала с ошибками
	go func() {
		for err := range chError {
			log.Println(err)
		}
	}()
	log.Println("Запуск сервера на порту :8081")
	err = http.ListenAndServe(":8081", srv.api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

// Парсинг отдельной rss ленты
func ParseUrls(url string, db *db.Interface, posts chan<- []db.Post, errors chan<- error, period int) {
	//log.Println("Внутри функции парсинга начало")
	ticker := time.NewTicker(time.Minute * time.Duration(period))
	defer ticker.Stop()
	for range ticker.C {
		news, err := rss.ParseRss(url)
		if err != nil {
			errors <- err
			continue
		}
		posts <- news
	}
	//	log.Println("Внутри функции парсинга конец")
}
