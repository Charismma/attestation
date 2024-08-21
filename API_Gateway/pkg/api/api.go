package api

import (
	"API_Gateway/pkg/models"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

const req_id_str = "request_id"

// Структура API
type API struct {
	router *mux.Router
}

// Конструктор API
func New() *API {
	api := API{
		router: mux.NewRouter(),
	}
	api.endpoints()
	return &api
}

// Регистрация обработчиков
func (api *API) endpoints() {
	api.router.Use(identMiddleware)
	api.router.Use(loggerMiddleware)
	api.router.HandleFunc("/news", api.news).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/news/filter", api.newsFilter).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/news/id", api.detailedNews).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/news/comment", api.addComment).Methods(http.MethodPost, http.MethodOptions)
}
func (api *API) Router() *mux.Router {
	return api.router
}

// Middleware логирования
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			req_id_string := r.Context().Value(req_id_str).(string)
			ip := r.RemoteAddr
			log.Println("Пришел новый запрос: request_id=", req_id_string, ", ip-адрес: ", ip)
		}()

		next.ServeHTTP(w, r)
	})
}

// Middleware для работы с идетификатором запроса
func identMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request_id := r.URL.Query().Get("request_id")
		if request_id == "" {
			//log.Println("Генерируем значение")
			request_id = generator_reqid()

		}
		//log.Println(request_id)
		ctx := context.WithValue(context.Background(), req_id_str, request_id)
		//log.Println("Запрос прошел через промежуточное ПО ident")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Обработчик получения новостей на странице
func (api *API) news(w http.ResponseWriter, r *http.Request) {
	req_id_string := r.Context().Value(req_id_str).(string)
	//log.Println("Запустился обработчик news")
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	urlStr := "http://localhost:8081/news?page=" + page + "&request_id=" + req_id_string
	resp, err := http.Get(urlStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)

}

// Обработчик получения отфильтрованных по строке новостей
func (api *API) newsFilter(w http.ResponseWriter, r *http.Request) {
	req_id_string := r.Context().Value(req_id_str).(string)
	//log.Println("Запустился обработчик newsfilter")
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	str := r.URL.Query().Get("s")
	urlStr := "http://localhost:8081/news/filter?s=" + str + "&page=" + page + "&request_id=" + req_id_string
	resp, err := http.Get(urlStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)

}

// Получение детализированной новости (Новости и комментарии к нему)
func (api *API) detailedNews(w http.ResponseWriter, r *http.Request) {
	req_ids := r.Context().Value(req_id_str).(string)
	//log.Println("Запустился обработчик detailedNews")
	id_news := r.URL.Query().Get("id")
	if id_news == "" {
		id_news = "1"
	}
	result := make(chan interface{}, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	var news []models.NewsFullDetailed
	var comment []models.Comment
	go func() {
		for res := range result {
			switch res := res.(type) {
			case error:
				//log.Println("Ошибка")
				wg.Done()
				http.Error(w, res.Error(), http.StatusBadRequest)
				return
			case []models.NewsFullDetailed:
				//log.Println(res)
				news = res
				wg.Done()
			case []models.Comment:
				//log.Println(res)
				comment = res
				wg.Done()
			}
		}
	}()
	go getNews(id_news, result, req_ids)
	go getComments(id_news, result, req_ids)
	wg.Wait()
	close(result)
	if len(news) != 0 {
		var newsWithComment []models.NewsFullDetailed
		newsWithComment = news
		newsWithComment[0].Comments = comment

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newsWithComment)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

// Обработчик создания нового комментария с проверкой по цензуре
func (api *API) addComment(w http.ResponseWriter, r *http.Request) {
	req_id_string := r.Context().Value(req_id_str).(string)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cenz_req, err := http.NewRequest("POST", "http://localhost:8083/cenz?request_id="+req_id_string, bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cenz_req.Header.Set("Content-Type", "application/json")
	cenz_client := &http.Client{}
	resp_cenz, err := cenz_client.Do(cenz_req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if resp_cenz.StatusCode == http.StatusOK {
		req, err := http.NewRequest("POST", "http://localhost:8082/addComment?request_id="+req_id_string, bytes.NewBuffer(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

// Получение новости по id
func getNews(id string, result chan<- interface{}, req_ids string) {
	urlStr := "http://localhost:8081/news/id?id=" + id + "&request_id=" + req_ids
	resp, err := http.Get(urlStr)
	if err != nil {
		result <- err
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result <- err
		return
	}
	news := []models.NewsFullDetailed{}
	err = json.Unmarshal(body, &news)
	if err != nil {
		result <- err
	}
	result <- news
}
func getComments(id string, result chan<- interface{}, req_ids string) {
	urlStr := "http://localhost:8082/comments?id_news=" + id + "&request_id=" + req_ids
	resp, err := http.Get(urlStr)
	if err != nil {
		result <- err
		return
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		result <- err
		return
	}
	comments := []models.Comment{}
	err = json.Unmarshal(body, &comments)
	if err != nil {
		result <- err
	}
	//log.Println("Десериализация и отправка по каналу", varunm)
	result <- comments
}

// генератор идентификатора запроса
func generator_reqid() string {
	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	var result string
	for i := 0; i < 10; i++ {
		j := rand.Intn(10)
		result += numbers[j]
	}
	return result
}
