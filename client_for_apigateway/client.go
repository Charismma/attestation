package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Comment struct {
	ID        int       `json:"-"`
	Post_id   int       `json:"Post_id"`
	Parent_id int       `json:"Parent_id"`
	Content   string    `json:"Content"`
	AddTime   int64     `json:"AddTime"`
	Visible   bool      `json:"-"`
	Replies   []Comment `json:"-"`
}

func main() {
	comm := &Comment{
		Post_id:   101,
		Parent_id: 0,
		Content:   "йцуке",
		AddTime:   0,
	}
	data, err := json.Marshal(comm)
	if err != nil {
		fmt.Println(err)
		return
	}
	req, err := http.NewRequest("POST", "http://localhost:8080/news/comment", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Выводим ответ от сервера
	fmt.Println(resp.Status)
}
