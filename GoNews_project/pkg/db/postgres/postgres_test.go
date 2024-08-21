package postgres

import (
	"GoNews_project/pkg/db"
	"strconv"
	"testing"
	"time"
)

var connstr = "postgres://postgres:password@192.168.1.191:5432/GoNews"

func TestStorage_Posts(t *testing.T) {
	var posts = []db.Post{
		{
			ID:      1,
			Title:   "Заголовк 1",
			Content: "Содержание 1",
			PubTime: time.Now().Unix(),
			Link:    strconv.Itoa(int(time.Now().Unix() + 50)),
		},
		{
			ID:      2,
			Title:   "Заголовк 2",
			Content: "Содержание 2",
			PubTime: time.Now().Unix() + 1,
			Link:    strconv.Itoa(int(time.Now().Unix() + 100)),
		},
	}
	db, err := New(connstr)
	if err != nil {
		t.Fatal(err)
	}
	err = db.AddPosts(posts)
	if err != nil {
		t.Log(err)
	}
	postss, err := db.Posts(2, 2)
	if err != nil {
		t.Log(err)
	}
	t.Log(postss)
}

func TestNew(t *testing.T) {

	_, err := New(connstr)
	if err != nil {
		t.Fatal(err)
	}

}
