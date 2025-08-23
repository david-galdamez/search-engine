package models

type Document struct {
	Id      string  `json:"id"`
	Title   string  `json:"title"`
	Content string  `json:"content"`
	Url     *string `json:"url,omitempty"`
	Length  int     `json:"length"`
}
