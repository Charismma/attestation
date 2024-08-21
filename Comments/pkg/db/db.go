package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
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

// подключение к бд
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

// Добавление комментария в БД
func (s *Storage) AddComment(comment Comment) error {
	_, err := s.db.Exec(context.Background(), `INSERT INTO comments_post(post_id,parent_id,content,addTime)
	VALUES($1,$2,$3,$4)`,
		&comment.Post_id,
		&comment.Parent_id,
		&comment.Content,
		&comment.AddTime,
	)
	if err != nil {
		return err
	}
	return nil
}

// Получение комментариев по id поста
func (s *Storage) Comments(post_id int) ([]Comment, error) {
	rows, err := s.db.Query(context.Background(), `SELECT id,post_id,parent_id,content,addTime FROM comments_post WHERE post_id=$1 ORDER BY addTime DESC`,
		post_id)
	if err != nil {
		return nil, err
	}
	var comments []Comment
	for rows.Next() {
		var comment Comment
		err = rows.Scan(
			&comment.ID,
			&comment.Post_id,
			&comment.Parent_id,
			&comment.Content,
			&comment.AddTime,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	rows.Close()
	//log.Println(comments)
	result_comm := transform_comment(comments)
	//log.Println("Другое:", ret_comm)
	return result_comm, nil
}

// Преобразование комментариев, где в структуре сразу есть комментарии, которые являются ответом на другие комментарии
func transform_comment(comm []Comment) []Comment {
	var endComment []Comment
	for _, comment := range comm {
		if comment.Parent_id == 0 {
			endComment = append(endComment, comment)

		}
	}
	for _, comment := range comm {
		if comment.Parent_id != 0 {
			for i, comm2 := range endComment {
				if comm2.ID == comment.Parent_id {
					endComment[i].Replies = append(endComment[i].Replies, comment)
				}
			}
		}
	}
	return endComment
}
