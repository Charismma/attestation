package models

type NewsFullDetailed struct {
	ID       int64     `json:"ID"`
	Title    string    `json:"Title"`
	Content  string    `json:"Content"`
	PubTime  int64     `json:"PubTime"`
	Link     string    `json:"Link"`
	Comments []Comment `json:"Comments"`
}

type NewsShortDetailed struct {
	ID      int
	Title   string
	PubTime int64
	Link    string
}

type Comment struct {
	ID        int       `json:"ID"`
	Post_id   int       `json:"Post_ID"`
	Parent_id int       `json:"Parent_ID"`
	Content   string    `json:"Content"`
	AddTime   int64     `json:"AddTime"`
	Visible   bool      `json:"Visible"`
	Replies   []Comment `json:"Replies"`
}
