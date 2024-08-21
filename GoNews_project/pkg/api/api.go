package api

import (
	"GoNews_project/pkg/db"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// структура нашего API
type API struct {
	db     db.Interface
	router *mux.Router
}

// структура с объектом пагинации и новостями,которые отдаются пользователю
type Paginats struct {
	News []db.Post
	Pag  Pag
}

// Объект пагинации
type Pag struct {
	N      int //номер текущей страницы
	Pages  int //количство страниц
	OnPage int //количество новостей на странице
}

const OnPage = 10

// Конструктор API
func New(db db.Interface) *API {
	api := API{
		db: db,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация обработчиков
func (api *API) endpoints() {
	api.router.Use(loggerMiddleware)
	api.router.HandleFunc("/news", api.postsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/news/id", api.postById).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/news/filter", api.filternews).Methods(http.MethodGet, http.MethodOptions)
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

// Обработчик Get-запроса на получение новости по id
func (api *API) postById(w http.ResponseWriter, r *http.Request) {
	id_str := r.URL.Query().Get("id")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	post, err := api.db.PostById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(post)

}

// Обработчик Get-запроса на получение новостей постранично
func (api *API) postsHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("page")
	n, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	count, err := api.db.Count()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//log.Println("Количество постов:", count)
	start := OnPage * (n - 1)
	pages := (count / OnPage) + 1
	post, err := api.db.Posts(start, OnPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pag := Pag{
		N:      n,
		Pages:  pages,
		OnPage: OnPage,
	}
	send := Paginats{
		News: post,
		Pag:  pag,
	}
	json.NewEncoder(w).Encode(send)

}

// Обработчик Get-запроса на получение новостей по искомой строке
func (api *API) filternews(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("page")
	n, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	str := r.URL.Query().Get("s")
	count, err := api.db.CountOfFilter(str)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	start := OnPage * (n - 1)
	post, err := api.db.Filter(str, start, OnPage)
	pages := (count / OnPage) + 1
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pag := Pag{
		N:      n,
		Pages:  pages,
		OnPage: OnPage,
	}
	send := Paginats{
		News: post,
		Pag:  pag,
	}
	json.NewEncoder(w).Encode(send)

}
