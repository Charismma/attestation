package rss

import (
	"GoNews_project/pkg/db"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
)

// Xml-структура rss-ленты
type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Chanel  struct {
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Link        string `xml:"link"`
		Items       []struct {
			Title       string `xml:"title"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
			Link        string `xml:"link"`
		} `xml:"item"`
	} `xml:"channel"`
}

// Парсинг rss-ленты
func ParseRss(url string) ([]db.Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var f Feed
	err = xml.Unmarshal(b, &f)
	if err != nil {
		return nil, err
	}
	var items []db.Post
	for _, item := range f.Chanel.Items {
		var post db.Post
		post.Title = item.Title
		post.Content = item.Description
		post.Content = strip.StripTags(post.Content)
		post.Link = item.Link
		item.PubDate = strings.ReplaceAll(item.PubDate, ",", "")
		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", item.PubDate)
		}
		if err == nil {
			post.PubTime = t.Unix()
		}
		items = append(items, post)
	}
	return items, nil
}
