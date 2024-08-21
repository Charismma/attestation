package api

import (
	"Comments/pkg/db"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// структура API
type API struct {
	db     *db.Storage
	router *mux.Router
}

// функция конструктор
func New(db *db.Storage) *API {
	api := API{
		db: db,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// регистрация обработчиков
func (api *API) endpoints() {
	api.router.Use(loggerMiddleware)
	api.router.HandleFunc("/addComment", api.addComments).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/comments", api.comments).Methods(http.MethodGet, http.MethodOptions)
}

func (api *API) Router() *mux.Router {
	return api.router
}

// Middleware для логирования
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			request_id := r.URL.Query().Get("request_id")
			ip := r.RemoteAddr
			log.Println("Пришел новый запрос: request_id=", request_id, ", ip-адрес: ", ip)

		}()

		next.ServeHTTP(w, r)
	})
}

// Обработчик добавления комментария
func (api *API) addComments(w http.ResponseWriter, r *http.Request) {

	var comm db.Comment
	err := json.NewDecoder(r.Body).Decode(&comm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	err = api.db.AddComment(comm)
	if err != nil {
		//log.Println("Тут возникла ошибка")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// обработчик получения комментариев по ID новости
func (api *API) comments(w http.ResponseWriter, r *http.Request) {
	id_str := r.URL.Query().Get("id_news")
	id_news, err := strconv.Atoi(id_str)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	comments, err := api.db.Comments(id_news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)

}
