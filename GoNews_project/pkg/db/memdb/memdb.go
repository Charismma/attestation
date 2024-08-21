package memdb

import "GoNews_project/pkg/db"

type Store struct {
}

//Функция конструктор
func New() (*Store, error) {
	return new(Store), nil
}

//Вывод всех постов
func (s *Store) Posts(n int) ([]db.Post, error) {
	return posts, nil
}

//Добавление постов
func (s *Store) AddPosts(posts []db.Post) error {
	return nil
}

var posts = []db.Post{
	{
		ID:      1,
		Title:   "Заголовк 1",
		Content: "Содержание 1",
		PubTime: 0,
		Link:    "Ссылка 1",
	},
	{
		ID:      2,
		Title:   "Заголовк 2",
		Content: "Содержание 2",
		PubTime: 0,
		Link:    "Ссылка 2",
	},
}
