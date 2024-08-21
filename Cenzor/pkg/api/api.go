package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type API struct {
	router *mux.Router
}

// Структура комментария
type Comment struct {
	ID        int       `json:"-"`
	Post_id   int       `json:"Post_id"`
	Parent_id int       `json:"Parent_id"`
	Content   string    `json:"Content"`
	AddTime   int64     `json:"AddTime"`
	Replies   []Comment `json:"-"`
}

// конструтор апи
func New() *API {
	api := API{
		router: mux.NewRouter(),
	}
	api.endpoints()
	return &api
}

// регистрация обработчиков
func (api *API) endpoints() {
	api.router.Use(loggerMiddleware)
	api.router.HandleFunc("/cenz", api.cenz).Methods(http.MethodPost, http.MethodOptions)
}

func (api *API) Router() *mux.Router {
	return api.router
}

// middleware логирования
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

// Обработчик для цензурирования комментариев
func (api *API) cenz(w http.ResponseWriter, r *http.Request) {
	var comm Comment
	err := json.NewDecoder(r.Body).Decode(&comm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	resp := cenzor(comm.Content)
	if resp {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}

// функция проверки тела комментария с имющимся списком слов
func cenzor(comm string) bool {
	var swear [3]string = [3]string{"йцукен", "ячсмит", "пролсд"}
	for i := 0; i < len(swear); i++ {
		if strings.Contains(comm, swear[i]) {
			return false
		}
	}
	return true
}
