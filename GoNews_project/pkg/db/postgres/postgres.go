package postgres

import (
	"GoNews_project/pkg/db"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

// Функция конструктор
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Вывод n постов со смещением и определенным количеством
func (s *Storage) Posts(start int, numOfPage int) ([]db.Post, error) {
	rows, err := s.db.Query(context.Background(), `SELECT id,title,content,pubtime,link FROM posts ORDER BY pubtime DESC OFFSET $1 LIMIT $2`, start, numOfPage)
	if err != nil {
		return nil, err
	}
	var posts []db.Post
	for rows.Next() {
		var post db.Post
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.PubTime,
			&post.Link,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()

}

// счетчик количества строк в таблице posts
func (s *Storage) Count() (int, error) {
	rows, err := s.db.Query(context.Background(), `SELECT COUNT(*) AS total_rows FROM posts`)
	if err != nil {
		return 0, err
	}
	var count int
	for rows.Next() {
		err = rows.Scan(
			&count,
		)
		if err != nil {
			return 0, err
		}
	}
	return count, rows.Err()

}

// получение детальной новости по id
func (s *Storage) PostById(n int) ([]db.Post, error) {
	rows, err := s.db.Query(context.Background(), `SELECT id,title,content,pubtime,link FROM posts WHERE id=$1`, n)
	if err != nil {
		return nil, err
	}
	var posts []db.Post
	for rows.Next() {
		var post db.Post
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.PubTime,
			&post.Link,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()

}

// получение отфилтьрованных новостей по заголовку со смещением и ограничением
func (s *Storage) Filter(str string, n int, numOfPage int) ([]db.Post, error) {
	str = "%" + str + "%"
	rows, err := s.db.Query(context.Background(), `SELECT id,title,content,pubtime,link FROM posts WHERE title ILIKE $1 ORDER BY pubtime DESC OFFSET $2 LIMIT $3`, str, n, numOfPage)
	if err != nil {
		return nil, err
	}
	var posts []db.Post
	for rows.Next() {
		var post db.Post
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.PubTime,
			&post.Link,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()

}

// счетчик количества новостей по фильтру
func (s *Storage) CountOfFilter(str string) (int, error) {
	str = "%" + str + "%"
	rows, err := s.db.Query(context.Background(), `SELECT COUNT(*) AS total FROM posts WHERE title ILIKE $1`, str)
	if err != nil {
		return 0, err
	}
	var count int
	for rows.Next() {
		err = rows.Scan(
			&count,
		)
		if err != nil {
			return 0, err
		}
	}
	return count, rows.Err()

}

// Добавление постов
func (s *Storage) AddPosts(posts []db.Post) error {
	for _, post := range posts {
		_, err := s.db.Exec(context.Background(), `INSERT INTO posts(title,content,pubtime,link) VALUES($1,$2,$3,$4)`,
			&post.Title,
			&post.Content,
			&post.PubTime,
			&post.Link,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
