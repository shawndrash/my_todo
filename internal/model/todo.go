package model

type Todo struct {
	Id        int64  `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}
