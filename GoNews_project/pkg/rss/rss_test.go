package rss

import (
	"testing"
)

func TestParseRss(t *testing.T) {
	posts, err := ParseRss("https://habr.com/ru/rss/hub/go/all/?fl=ru")
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) == 0 {
		t.Fatal("Не получено постов")
	}
	t.Log("Количество новостей:", len(posts), " Посты:", posts)
}
